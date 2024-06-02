package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {

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

	//we dont want to apply sessions to serve static files
	fileServer := http.FileServer(http.Dir("./ui/static"))
	//update pattern for static resource path
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	//custom error for httprouter
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	//get on create returns a snippet form
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	//post request for creating snippet handled by this
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))
	//alice package makes middleware chaining easier
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	return standard.Then(router)
}
