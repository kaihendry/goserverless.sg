package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/tj/go/http/response"
)

func init() {
	if os.Getenv("UP_STAGE") == "" {
		log.SetHandler(text.Default)
	} else {
		log.SetHandler(jsonhandler.Default)
	}
}

func main() {
	addr := ":" + os.Getenv("PORT")
	app := mux.NewRouter()

	app.HandleFunc("/", handlePost).Methods("POST")
	app.HandleFunc("/rank", handleRank).Methods("GET")
	app.HandleFunc("/", handleIndex).Methods("GET")

	var options []csrf.Option
	// If developing locally
	if os.Getenv("UP_STAGE") == "" {
		// https://godoc.org/github.com/gorilla/csrf#Secure
		log.Warn("CSRF insecure")
		options = append(options, csrf.Secure(false))
	}

	if err := http.ListenAndServe(addr,
		// Only protects the POST btw
		csrf.Protect([]byte("go-serverless"), options...)(app)); err != nil {
		log.WithError(err).Fatal("error listening")
	}
}

func countryFlag(x string) string {
	if len(x) != 2 {
		return ""
	}
	if x[0] < 'A' || x[0] > 'Z' || x[1] < 'A' || x[1] > 'Z' {
		return ""
	}
	return string(0x1F1E6+rune(x[0])-'A') + string(0x1F1E6+rune(x[1])-'A')
}

func handleIndex(w http.ResponseWriter, r *http.Request) {

	// Don't index beta.goserverless.sg
	if os.Getenv("UP_STAGE") != "production" {
		w.Header().Set("X-Robots-Tag", "none")
	}

	t := template.Must(template.New("").ParseGlob("templates/*.html"))
	err := t.ExecuteTemplate(w, "index.html", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"Stage":          os.Getenv("UP_STAGE"),
		"Year":           time.Now().Format("2006"),
		"EmojiCountry":   countryFlag(strings.Trim(r.Header.Get("Cloudfront-Viewer-Country"), "[]")),
	})

	if err != nil {
		log.WithError(err).Error("template failed to parse")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func handlePost(w http.ResponseWriter, r *http.Request) {

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.WithError(err).Error("dumping request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(0)
	if err != nil {
		log.WithError(err).Error("parsing form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for key, values := range r.PostForm { // range over map
		for _, value := range values { // range over []string
			log.Infof("Key: %v Value: %v", key, value)
		}
	}

	cfg, err := external.LoadDefaultAWSConfig(external.WithSharedConfigProfile("gosls"))
	if err != nil {
		log.WithError(err).Error("loading config")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cfg.Region = endpoints.UsWest2RegionID

	svc := ses.New(cfg)
	input := &ses.SendEmailInput{
		Source: aws.String(fmt.Sprintf("%s <hendry@goserverless.sg>", r.PostFormValue("name"))),
		Destination: &ses.Destination{
			ToAddresses: []string{
				"hendry@goserverless.sg",
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(string(dump)),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(fmt.Sprintf("Inquiry from %s", r.PostFormValue("organization"))),
			},
		},
	}

	req := svc.SendEmailRequest(input)
	result, err := req.Send(context.TODO())
	if err != nil {
		log.WithError(err).Error("sending mail")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.JSON(w, result)
}

func handleRank(w http.ResponseWriter, r *http.Request) {

	type AWSRegion struct {
		Name         string `json:"name"`
		ServiceCount int    `json:"count"`
	}
	var regions []AWSRegion

	resolver := endpoints.NewDefaultResolver()
	partitions := resolver.Partitions()

	for _, p := range partitions {
		for id, r := range p.Regions() {
			services := r.Services()
			if id == "ap-southeast-1" {
				log.Infof("Service count in Singapore: %d", len(services))
			}
			regions = append(regions, AWSRegion{id, len(services)})
		}
	}
	sort.Slice(regions, func(i, j int) bool {
		return regions[i].ServiceCount > regions[j].ServiceCount
	})
	response.JSON(w, regions)
}
