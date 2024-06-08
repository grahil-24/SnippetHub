package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form"
	"net/http"
	"runtime/debug"
)

func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request, in the same way that we did in our
	// createSnippetPost handler.
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		// For all other errors, we return them as normal.
		return err
	}
	return nil
}

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

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {

	ts, ok := app.templateCache[page]
	if !ok {
		app.serverError(w, fmt.Errorf("page %s not found in template cache", page))
		return
	}

	//initialize a new buffer
	buf := new(bytes.Buffer)

	//using buffer we can catch runtime errors in rendering the templates, so users wont get
	//a 200 status code even if the page was not rendered properly. We first render the template
	// in a buffer. If it does not return an error, we write to the ResponseWriter from the buffer
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)

}

// Return true if the current request is from an authenticate user, otherwise
// return false
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
