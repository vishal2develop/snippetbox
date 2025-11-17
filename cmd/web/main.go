package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.vishalborana2407.net/internal/models"
)

// _ = Import this package only for its side effects, not because I’m directly using its functions or types.

// Define an application struct to hold the application-wide dependencies
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel // will allow us to use the SnippetModel type in our handlers.
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	// accept a new command line flag for the port
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command-line flag for the MySQL DSN string.
	// web = username, admin = password, snippetbox = database name, parseTime = true = parse time
	dsn := flag.String("dsn", "web:web@tcp(127.0.0.1:3306)/snippetbox?parseTime=true", "MySQL DSN string")

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

	// Initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// initialize a new form decoder
	formDecoder := form.NewDecoder()

	// initialize a new session manager, then configure it to use our MySQL database as our session store.
	// set lifetime to 12 hours, so that the session cookie expires automatically after 12 hours.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Initialize a new instance of our application struct, containing the
	// dependencies
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db}, // contains the connection pool
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
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
