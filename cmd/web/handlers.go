package main

import (
	"errors"
	"fmt"
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

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)

	data.Snippets = snippets

	// log length of snippets
	app.logger.Info("Number of snippets", "length", len(snippets))

	// use the render helper to render the home.tmpl template
	app.render(w, r, http.StatusOK, "home.tmpl", data)
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

	data := app.newTemplateData(r)

	data.Snippet = snippet

	// Use the render helper.
	app.render(w, r, http.StatusOK, "view.tmpl", data)
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
