package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// The serverError helper writes a log entry at Error level (including the request
// method and URI as attributes), then sends a generic 500 Internal Server Error
// response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = debug.Stack()
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri, "stack_trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Retrieve the appropriate template set from the cache based on the page
	// name (like 'home.tmpl'). If no entry exists in the cache with the
	// provided name, then create a new error and call the serverError() helper and return

	// get the template from the cache
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// Phase 1
	// initialize a new buffer to store the rendered template
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError() helper
	// and then return.

	err := ts.ExecuteTemplate(buf, "base", data)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Phase 2

	// If the template is written to the buffer without any errors, it's safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter.
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

// newTemplateData creates a new templateData struct intialized with the current year.
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}

// helper utiity for form parsing + decoding and checking for errors
// Create a new decodePostForm() helper method. The second parameter here, dst,
// is the target destination into which we want to decode the form data.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// parseForm
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// decode using form decoder
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		//Handle InvalidDecoderError separately
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type form.InvalidDecoderError. We use
		// errors.As() to check for this and panic.
		// InvalidDecoderError occurs when the target destination is not a non-nil pointer.
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// for all other errors, return the error
		return err
	}
	return nil
}
