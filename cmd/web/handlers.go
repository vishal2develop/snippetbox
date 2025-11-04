package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"snippetbox.vishalborana2407.net/internal/models"
)

// home handles requests to the root URL ("/").
// Change the signature of the home handler so it is defined as a method against
// *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Add a custom response header for demonstration.
	w.Header().Add("Server", "Go")

	// Get latest snippet - top 10
	snippets, err := app.snippets.Latest()

	// If there’s an error in getting the records, return server error - 500
	if err != nil {
		app.serverError(w, r, err)
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Create an instance of a templateData struct holding the slice of
	// snippets.
	data := templateData{
		Snippets: snippets,
	}

	// Pass in the templateData struct when executing the template.
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}

}

// snippetView handles requests for viewing a specific snippet.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// ✅ Extract the "id" parameter from the URL and convert it to an integer.
	// 💡 PathValue() gets the dynamic value from the route pattern, e.g. /snippet/view/{id}
	id, err := strconv.Atoi(r.PathValue("id"))

	// ⚠️ If the ID is invalid (non-numeric or less than 1), return a 404 page.
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Initialize a slice containing the paths to the view.tmpl file,
	// plus the base layout and navigation partial that we made earlier.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	// Parse and Create a new template from the files.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Create a new template data struct to pass to the template.
	data := templateData{Snippet: snippet}

	// Execute the template with the snippet data.
	err = ts.ExecuteTemplate(w, "base", data)

	if err != nil {
		app.serverError(w, r, err)
	}

	// Write the snippet data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%+v", snippet)
}

// snippetCreate displays a form for creating a new snippet.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// snippetCreatePost handles form submissions and saves a new snippet.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// dummy data for snippet creation
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// call the Insert() method on the snippet model
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
