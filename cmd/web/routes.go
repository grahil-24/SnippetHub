package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {

	fileServer := http.FileServer(http.Dir("./ui/static"))
	//
	//mux := http.NewServeMux()
	//mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	//mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/snippet/view", app.snippetView)
	//mux.HandleFunc("/snippet/create", app.snippetCreate)
	//
	////pass the serveMux as the next parameter to our middleware, so as to use this middleware for every incoming request
	////is used to apply secure response headers
	////flow of middlewares - logRequest -> secureHeader -> ServeMux -> ApplicationHandler
	////in panic recovery, the middleware is executed in only the goroutine in which it was executed. If another
	////goroutine, was spun from it, it wont happen in that
	////return app.recoverPanic(app.logRequest(secureHeader(mux)))
	//

	//third party router provide routing through requets method too
	router := httprouter.New()
	//custom error for httprouter
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	//update pattern for static resource path
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	//get on create returns a snippet form
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	//post request for creating snippet handled by this
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)
	//alice package makes middleware chaining easier
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	return standard.Then(router)
}
