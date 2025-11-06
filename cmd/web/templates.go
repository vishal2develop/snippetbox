package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.vishalborana2407.net/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the project progresses.

// lowercase starting  = private (not accessible outside of this package)
// uppercase starting = public (accessible outside of this package)
type templateData struct {
	Snippet     models.Snippet
	Snippets    []models.Snippet
	CurrentYear int
}

// helper function to format a time.Time object as a human-readable date
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// initialize a template.Funcmap value and store it in a global variabe
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// create a new template cache that will hold all the templates
// Returns a map of template name to template object
func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl".
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")

	if err != nil {
		return nil, err
	}

	// build a full template set for each page individually.
	for _, page := range pages {
		// Extract the filename from the filepath.
		name := filepath.Base(page)

		// before we parse the template, register the functions
		// template.New(name) - Creates a new, empty template with the given name
		// .Funcs(functions) - Registers custom functions (like humanDate) that can be used in templates
		ts := template.New(name).Funcs(functions)

		// Parse the base template file into a template set.
		// Start with the master layout template.
		ts, err := ts.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		// Add all partial templates (header, footer, etc.) to the same template set.
		// Now ts knows about: base.tmpl, nav.tmpl and future partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the page template.
		// Parse the page template itself
		// Now ts knows about: base.tmpl, partials and the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template to the cache.
		// key = template name, value = template object
		cache[name] = ts
	}
	// Return the map.
	return cache, nil
}
