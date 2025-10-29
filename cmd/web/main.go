package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log/slog"
	"net/http"
	"os"
)

// _ = Import this package only for its side effects, not because I’m directly using its functions or types.

// Define an application struct to hold the application-wide dependencies
type application struct {
	logger *slog.Logger
}

func main() {

	// accept a new command line flag for the port
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command-line flag for the MySQL DSN string.
	// web = username, admin = password, snippetbox = database name, parseTime = true = parse time
	dsn := flag.String("dsn", "web:admin@/snippetbox?parseTime=true", "MySQL DSN string")

	// parse the flags and assign it to addr.
	// Parse() must be called after all flags are defined and before flags are accessed.
	// if not called, the flag will be set to the default value.
	// any errors during parsing, the application will be terminated.
	flag.Parse()

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	// second argument is a pointer to a slog.HandlerOptions struct , which you can use to customize the behavior of the handler. if happy, with default settings -> pass nil
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger).
	app := &application{
		logger: logger,
	}

	// Value returned by flag.String() is a pointer to the flag's value and not the value itself.
	// Hence, we need to dereference the pointer (prefix with *) to get the actual value.
	logger.Info("Starting server on", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	// terminate the application with exit code 1.
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
