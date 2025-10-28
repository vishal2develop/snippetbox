package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies
type application struct {
	logger *slog.Logger
}

func main() {

	// accept a new command line flag for the port
	addr := flag.String("addr", ":4000", "HTTP network address")

	// parse the flags and assign it to addr.
	// Parse() must be called after all flags are defined and before flags are accessed.
	// if not called, the flag will be set to the default value.
	// any errors during parsing, the application will be terminated.
	flag.Parse()

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	// second argument is a pointer to a slog.HandlerOptions struct , which you can use to customize the behavior of the handler. if happy, with default settings -> pass nil
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger).
	app := &application{
		logger: logger,
	}

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
	// Swap the route declarations to use the application struct's methods as the
	// handler functions.
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	// Register POST routes
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Value returned by flag.String() is a pointer to the flag's value and not the value itself.
	// Hence, we need to dereference the pointer (prefix with *) to get the actual value.
	logger.Info("Starting server on", "addr", *addr)

	err := http.ListenAndServe(*addr, mux)

	logger.Error(err.Error())
	// terminate the application with exit code 1.
	os.Exit(1)
}
