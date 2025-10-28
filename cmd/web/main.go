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

	// Value returned by flag.String() is a pointer to the flag's value and not the value itself.
	// Hence, we need to dereference the pointer (prefix with *) to get the actual value.
	logger.Info("Starting server on", "addr", *addr)

	err := http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	// terminate the application with exit code 1.
	os.Exit(1)
}
