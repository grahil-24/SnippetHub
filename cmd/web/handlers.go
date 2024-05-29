package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// "/" route handler. home page fetches the top 10 latest snippets
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	//get PWD
	//wd, err := os.Getwd()
	//fmt.Println(wd)

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

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	app.render(w, http.StatusOK, "create.gohtml", app.newTemplateData(r))
}

// "/snippet/create" route handler
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Print("inside snippetCreate")

	//checking if the method is post is no longer required and can be removed

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
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	//id, err := strconv.Atoi(r.URL.Query().Get("id"))

	//fetching parameters. returns a slice
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))

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
