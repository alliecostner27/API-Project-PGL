package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Sample data to pass to the template
	tags := map[string][]string{
		"Tags": {"dessert", "italian"},
	}

	// Serve static files (images)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	// Serve static HTML files directly
	http.Handle("/randomRecipe.html", http.FileServer(http.Dir(".")))

	// Handler for the homepage
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, tags)
	}

	// Handle form submission
	h2 := func(w http.ResponseWriter, r *http.Request) {
		var responseTags []string
		for key := range tags["Tags"] {
			if r.PostFormValue(tags["Tags"][key]) == "on" {
				responseTags = append(responseTags, tags["Tags"][key])
			}
		}
	
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h2>Hello</h2> <p>You selected: %v</p>", responseTags)
	}

	// Register handlers for routes
	http.HandleFunc("/", h1)
	http.HandleFunc("/send-tag/", h2)

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", nil))
}

