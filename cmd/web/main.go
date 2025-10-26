package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
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
