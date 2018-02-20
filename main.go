package main

import (
	"fmt"
	"net/http"
	"os"

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

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// signup_form.tmpl just needs a {{ .csrfField }} template tag for
	// csrf.TemplateField to inject the CSRF token into. Easy!
	t := template.Must(template.New("").ParseGlob("templates/*.html"))
	t.ExecuteTemplate(w, "index.html", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"Header":         r.Header,
	})
	// We could also retrieve the token directly from csrf.Token(r) and
	// set it in the request header - w.Header.Set("X-CSRF-Token", token)
	// This is useful if you're sending JSON to clients or a front-end JavaScript
	// framework.
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Oh hai")
}
