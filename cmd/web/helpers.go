package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	// The serverError helper writes an error message and stack trace to the errorLog,
	// then sends a generic 500 Internal Server Error response to the user.
	trace := fmt.Sprintf("%sn%s", err.Error(), debug.Stack())
	//setting call depth to 2 as at the top of stack its the helper file as the error is written here
	app.errorLog.Output(2, trace)
	http.Error(w, trace, http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// special error handler for 404s
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
