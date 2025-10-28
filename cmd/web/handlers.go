package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// home handles requests to the root URL ("/").
// Change the signature of the home handler so it is defined as a method against
// *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Add a custom response header for demonstration.
	w.Header().Add("Server", "Go")

	// ✅ List of template files to parse.
	// 💡 The base layout (base.tmpl) should always come first,
	// because it defines the common HTML structure (like <html>, <head>, etc.)
	// that other page templates (e.g. home.tmpl) will embed into.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}

	// ✅ Parse all template files into a single *template.Template object.
	// 💡 The "..." (ellipsis) expands the slice so that each file path
	// is passed as an individual argument to template.ParseFiles().
	templateSet, err := template.ParseFiles(files...)

	// 💡 Log all defined templates (useful for debugging).
	//log.Print("templateSet: ", templateSet.DefinedTemplates())

	// ⚠️ Always check the error immediately after parsing.
	// If there’s a syntax or path error, return a 500 response and stop further execution.
	if err != nil {
		// log at Error level containing the error message,
		//also including the request method and URI as attributes to assist with debugging.
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// ✅ Execute the base template and write the output to the ResponseWriter.
	// 💡 The "base" template usually includes placeholders ({{template "title"}} / {{template "main"}})
	// that will automatically call the nested templates like home.tmpl.
	// The second argument allows passing dynamic data to the template (nil here means no data).
	// why pass "base": Render and execute the template named base, and write the final HTML to w (the HTTP response).
	err = templateSet.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	// ✅ Write a response that displays the snippet ID.
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// snippetCreate displays a form for creating a new snippet.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// snippetCreatePost handles form submissions and saves a new snippet.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// ✅ Respond with a 201 Created status to indicate successful creation.
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}
