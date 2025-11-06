package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.vishalborana2407.net/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the project progresses.
type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
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

		// Parse the base template file into a template set.
		// Start with the master layout template.
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
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
