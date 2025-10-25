package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Component 1 (handler): Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
/**
* w http.ResponseWriter: An interface used to construct and send the HTTP response back to the client. It provides methods like Write(), WriteHeader(), and Header().
* r *http.Request: A pointer to a struct containing all information about the incoming HTTP request (URL, headers, method, body, etc.).
 */
func home(w http.ResponseWriter, r *http.Request) {
	// Log that the home handler was called
	log.Println("Home handler called")
	// []byte("Hello from Snippetbox"): Converts the string literal to a byte slice, which is the format Write() expects
	// w.Write(): Sends the byte slice as the HTTP response body to the client
	w.Write([]byte("Hello from Snippetbox"))
}

// handler to view snippet
func viewSnippet(w http.ResponseWriter, r *http.Request) {
	log.Println("View Snippet handler called")
	// Extract the value of the id wildcard from the request using r.PathValue()
	// and try to convert it to an integer using the strconv.Atoi() function. If
	// it can't be converted to an integer, or the value is less than 1, we
	// return a 404 page not found response.
	snippetId := r.PathValue("id")
	id, err := strconv.Atoi(snippetId)

	if err != nil || id < 1 {
		log.Println("Invalid snippet ID: ", snippetId)
		log.Println("Error: ", err)
		http.NotFound(w, r)
		return
	}
	// interpolate response string
	msg := fmt.Sprintf("Display the snippet with ID: %d", id)
	log.Println("Snippet ID: ", snippetId)
	w.Write([]byte(msg))
}

// handler to create snippet
func createSnippet(w http.ResponseWriter, r *http.Request) {
	log.Println("Create Snippet handler called")
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func main() {
	//Componenet 2 (Router/Servermux): Use the http.NewServeMux() function to initialize a new servemux (router), then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home) // {$} = Restrict this route to exact matches on / only.
	mux.HandleFunc("/snippet/view/{id}", viewSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Not Recommended
	// Register routes without explicitly declaring a servemux (http.handle/ http.handleFunc)
	// Internally uses the default servermux created automatically by GO, stored as a global variable (http.DefaultServeMux).
	// http.HandleFunc("/snippet/view", viewSnippet)

	// Print a log message to say that the server is starting.
	log.Print("starting server on :4000")

	// Component 3 (Server): Use the http.ListenAndServe() function to start a server listening on port 4000.
	// The server will listen for incoming HTTP requests and call the handler function registered for the URL pattern.
	err := http.ListenAndServe(":4000", mux)

	// If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and terminate the
	// program. Note that any error returned by http.ListenAndServe() is always non-nil.
	log.Fatal(err)

	// Not Recommended
	// if we pass nil as the second argument to http.ListenAndServe()
	// it will use the default servermux
	//err := http.ListenAndServe(":4000", nil)

}
