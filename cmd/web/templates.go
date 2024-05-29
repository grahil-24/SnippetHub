package main

import (
	"html/template"
	"path/filepath"
	"snippetbox.rahilganatra.net/internal/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplatesCache() (map[string]*template.Template, error) {
	//init a new map which acts as cache
	cache := map[string]*template.Template{}

	//returns slices of path which match this pattern
	pages, err := filepath.Glob("./ui/html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		//extract the file name from the path
		name := filepath.Base(page)

		files := []string{"./ui/html/pages/base.gohtml", "./ui/html/partials/nav.gohtml", page}

		t, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = t
	}
	return cache, nil
}
