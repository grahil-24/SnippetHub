package main

import (
	"fmt"
	"net/http"
)

//all custom middlewares  go here

//dealing with panic
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// want this middleware to be executed on every request
func secureHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//restricts where the resources like javascript, fonts, images, css can be loaded from
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		//sets if the page can be embedded in other sites. prevents clickjacking attacks
		w.Header().Set("X-Frame-Options", "DENY")
		//to prevent MIME sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		//0 is to turn off XSS. Use the better alternative ie the Content-Security-Policy
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		//code above this is executed down the chain
		//if the execution returns before next() is called, the flow is upstreamed, and we cant go further down the chain
		next.ServeHTTP(w, r)
		//code below next is executed up the chain
	})
}

//middleware to log every request's info like method, ip, protocol etc.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s-%s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}
