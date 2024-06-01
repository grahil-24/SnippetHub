package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"snippetbox.rahilganatra.net/internal/models"
	"time"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
}

// returns formatted time
func humanDate(t time.Time) string {
	return t.Format("Mon Jan _2 15:04:05 2006")
}

// global map for template functions. Functions need to be registered, before the template is parsed
// so that the functions can be used inside the templates. Template functions should always return 1 value
// not more than that. Or 1 value and 1 error
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
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

		//parse the base template into a template set and register the functions to the template before
		//ParseFiles() is called
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/pages/base.gohtml")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/nav.gohtml")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
