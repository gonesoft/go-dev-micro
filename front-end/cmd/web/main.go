package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Embed all .gohtml files in the templates directory.
// This removes any dependency on absolute or relative filesystem paths at runtime.
//
//go:embed templates/*.gohtml
var templateFS embed.FS

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET / -> rendering test.page.gohtml")
		render(w, "test.page.gohtml")
	})

	log.Println("Starting front end service on port 80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string) {
	tmpl, err := template.ParseFS(
		templateFS,
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
		fmt.Sprintf("templates/%s", t),
	)
	if err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, "template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// If the base template is defined with {{define "base"}}, render it explicitly.
	if tmpl.Lookup("base") != nil {
		if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
			log.Printf("template exec error (base): %v", err)
			http.Error(w, "template exec error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// If no named "base" template, try executing the parsed root.
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("template exec error (root): %v", err)
		http.Error(w, "template exec error: "+err.Error(), http.StatusInternalServerError)
	}
}
