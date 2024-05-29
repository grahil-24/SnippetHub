package main

import "net/http"

func (app *application) routes() http.Handler {

	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	//pass the serveMux as the next parameter to our middleware
	return secureHeader(mux)

}
