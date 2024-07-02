package main

import (
	"context"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

//all custom middlewares  go here

// dealing with panic
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

// middleware to log every request's info like method, ip, protocol etc.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s-%s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//if the request is not authenticate, means user is not logged in, we redirect the user
		//to the home page
		if !app.isAuthenticated(r) {
			//add the path user is trying to access, to their session data
			app.sessionManager.Put(r.Context(), "redirectPathAfterLogin", r.URL.Path)

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		//else we add this header, which means that pages requiring authentication are not
		//stored in cache
		w.Header().Add("Cache-Control", "no-store")

		//and call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly attributes set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the authenticatedUserID value from the session using the
		// GetInt() method. This will return the zero value for an int (0) if no
		// "authenticatedUserID" value is in the session -- in which case we
		// call the next handler in the chain as normal and return.
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")

		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Otherwise, we check to see if a user with that ID exists in our
		// database.
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}
		//call the next handler in chain
		next.ServeHTTP(w, r)
	})
}
