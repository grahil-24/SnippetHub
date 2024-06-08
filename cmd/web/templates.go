package main

import (
	"github.com/justinas/nosurf"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"snippetbox.rahilganatra.net/internal/models"
	"snippetbox.rahilganatra.net/ui"
	"time"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// returns formatted time
func humanDate(t time.Time) string {
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// global map for template functions. Functions need to be registered, before the template is parsed
// so that the functions can be used inside the templates. Template functions should always return 1 value
// not more than that. Or 1 value and 1 error
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	//fetch the flash message we stored in session during creation of snippet.
	//PopString() is sort of one time use. As we want to flash the message once after the snippet
	//has been created, we will fetch and remove it from the session. PopString() does exactly this
	//If there was no key matching with "flash" an empty string will be returned.
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func newTemplatesCache() (map[string]*template.Template, error) {
	//init a new map which acts as cache
	cache := map[string]*template.Template{}

	//returns slices of path which match this pattern. fetching the templates from the Embedded file system.
	pages, err := fs.Glob(ui.Files, "html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		//extract the file name from the path
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates we
		// want to parse.
		patterns := []string{
			"html/pages/base.gohtml",
			"html/partials/*.gohtml",
			page,
		}
		// Use ParseFS() instead of ParseFiles() to parse the template files
		// from the ui.Files embedded filesystem.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
