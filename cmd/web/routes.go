package main

import "net/http"
import "github.com/justinas/alice"

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	// Register static files
	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer((http.Dir("./ui/static")))

	// Registers a handler for any URL path that starts with /static/.
	/**
	Why Strip the Prefix?
	Without http.StripPrefix, there would be a mismatch:
	Incoming request: GET /static/css/style.css
	FileServer looks for: ./ui/static/static/css/style.css ❌ (looks for "static" twice!)
	*/
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// create a dynamic middleware chain that will be used by only specific routes and will contain the session middleware
	dynamic := alice.New(app.sessionManager.LoadAndSave)
	// Register handlers
	// Swap the route declarations to use the application struct's methods as the
	// handler functions.

	// Routes to use the dynamic middleware chain
	// handleFunc: wraps your function into an http.Handler. Expects: func(w, r)
	// handle: expects: http.Handler
	// ThenFunc: ThenFunc takes a normal func(w, r) and returns a proper http.Handler
	// Register GET routes
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))

	// Register POST routes
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// create standard middleware chain that will be used by all routes
	standardChain := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standardChain.Then(mux)
}
