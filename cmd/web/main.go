package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Register static files
	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Registers a handler for any URL path that starts with /static/.
	/**
	Why Strip the Prefix?
	Without http.StripPrefix, there would be a mismatch:
	Incoming request: GET /static/css/style.css
	FileServer looks for: ./ui/static/static/css/style.css ❌ (looks for "static" twice!)
	*/
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Register GET routes
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	// Register POST routes
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Println("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}
