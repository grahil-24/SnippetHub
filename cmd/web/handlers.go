package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

// "/" route handler
func home(w http.ResponseWriter, r *http.Request) {

	wd, err := os.Getwd()
	fmt.Println(wd)

	if r.URL.Path != "/" {
		//404 error(resource not found error)
		http.NotFound(w, r)
		return
	}
	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and use
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	files := []string{"./ui/html/pages/base.gohtml", "./ui/html/partials/nav.gohtml", "./ui/html/pages/home.gohtml"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
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
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

// "/snippet/create" route handler
func snippetCreate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {

		w.Header().Set("Allow", http.MethodPost)
		//only possible to call this method once per response
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a snippet"))

}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if id < 1 || err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "snippet #%d", id)

}
