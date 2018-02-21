package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"html/template"

	"github.com/apex/log"
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
	// options = append(options, csrf.Secure(false))

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
	fmt.Fprintf(w, "Oh hai")
}
