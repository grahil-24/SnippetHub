package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"snippetbox.rahilganatra.net/ui"
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

	//we dont want to apply sessions to serve static files. Getting the static files from the EFS
	fileServer := http.FileServer(http.FS(ui.Files))
	//update pattern for static resource path
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	//custom error for httprouter
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	//routers for user authentication and authorization
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	//these are protected routes, means these are only accessible to logged in users
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	//get on create returns a snippet form
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	//post request for creating snippet handled by this
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))

	//alice package makes middleware chaining easier
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	return standard.Then(router)
}
