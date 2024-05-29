package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// "/" route handler. home page fetches the top 10 latest snippets
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	wd, err := os.Getwd()
	fmt.Println(wd)

	if r.URL.Path != "/" {
		//404 error(resource not found error)
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	//for _, snippet := range snippets {
	//	fmt.Fprintf(w, "%+v", snippet)
	//}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.gohtml", data)

}

// "/snippet/create" route handler
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Print("inside snippetCreate")
	if r.Method != http.MethodPost {

		w.Header().Set("Allow", http.MethodPost)
		//only possible to call this method once per response
		w.WriteHeader(http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	//redirect to the created snippet after it has been created
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if id < 1 || err != nil {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	fmt.Println(snippet)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	//wrapping our data in a struct
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.gohtml", data)
}
