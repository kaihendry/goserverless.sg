package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"html/template"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/gorilla/csrf"
	"github.com/gorilla/pat"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	app := pat.New()

	app.Post("/", handlePost)
	app.Get("/", handleIndex)

	var options []csrf.Option
	// If developing locally
	options = append(options, csrf.Secure(false))

	if err := http.ListenAndServe(addr,
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
	t.ExecuteTemplate(w, "index.html", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"Stage":          os.Getenv("UP_STAGE"),
		"EmojiCountry":   countryFlag(strings.Trim(r.Header.Get("Cloudfront-Viewer-Country"), "[]")),
	})
}

func handlePost(w http.ResponseWriter, r *http.Request) {

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// var payload string
	// payload, err := ioutil.ReadAll(r.Body)
	// defer r.Body.Close()
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	err = r.ParseMultipartForm(0)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// fmt.Println(r.Form)
	// fmt.Println(r.PostForm)
	// fmt.Println(r.PostFormValue("organization"))
	// // fmt.Println(r.Body)

	// for key, values := range r.PostForm { // range over map
	// 	for _, value := range values { // range over []string
	// 		log.Infof("Key: %v Value: %v", key, value)
	// 	}
	// }

	cfg, err := external.LoadDefaultAWSConfig(external.WithSharedConfigProfile("gosls"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cfg.Region = endpoints.UsWest2RegionID

	svc := ses.New(cfg)
	input := &ses.SendEmailInput{
		Source: aws.String(fmt.Sprintf("%s %s <hendry@goserverless.sg>", r.PostFormValue("given-name"), r.PostFormValue("family-name"))),
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
	result, err := req.Send()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonb, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(jsonb))
}
