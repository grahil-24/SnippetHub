package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

// "/" route handler
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	wd, err := os.Getwd()
	fmt.Println(wd)

	if r.URL.Path != "/" {
		//404 error(resource not found error)
		app.notFound(w)
		return
	}
	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and use
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	files := []string{"./ui/html/pages/base.gohtml", "./ui/html/partials/nav.gohtml", "./ui/html/pages/home.gohtml"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// We then use the Execute() method on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	// leave as nil.
	//err = ts.Execute(w, nil)

	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.
	err = ts.ExecuteTemplate(w, "base", nil)

	if err != nil {
		app.serverError(w, err)
		return
	}
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
	fmt.Fprintf(w, "snippet #%d", id)

}
