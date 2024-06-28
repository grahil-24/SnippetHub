package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"snippetbox.rahilganatra.net/internal/models"
	"snippetbox.rahilganatra.net/internal/validator"
	"strconv"
)

// Update our snippetCreateForm struct to include struct tags which tell the
// decoder how to map HTML form values into the different struct fields. So, for
// example, here we're telling the decoder to store the value from the HTML form
// input with the name "title" in the Title field. The struct tag `form:"-"`
// tells the decoder to completely ignore a field during decoding
type snippetCreationForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// to check the status of the application
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "about.gohtml", data)
}

// "/" route handler. home page fetches the top 10 latest snippets
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	//get PWD
	//wd, err := os.Getwd()
	//fmt.Println(wd)

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	//for _, snippet := range snippets {
	//	fmt.Fprintf(w, "%+v", snippet)
	//}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.gohtml", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreationForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.gohtml", data)
}

// "/snippet/create" route handler
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Print("inside snippetCreate")

	var form snippetCreationForm

	// Call the Decode() method of the form decoder, passing in the current
	// request and *a pointer* to our snippetCreateForm struct. This will
	// essentially fill our struct with the relevant values from the HTML form.
	// If there is a problem, we return a 400 Bad Request response to the client.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Because the Validator type is embedded by the snippetCreateForm struct,
	// we can call CheckField() directly on it to execute our validation checks
	form.CheckField(validator.NotBlank(form.Title), "title", "this field cannot be empty")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "this field cannot be longer than 100 characters")
	form.CheckField(validator.PermittedValue(form.Expires, 365, 7, 1), "expires", "This field must equal 1, 7 or 365")
	form.CheckField(validator.NotBlank(form.Content), "content", "this field cannot be empty")

	//if error exists then return them in plain text http response
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.gohtml", data)
		return
	}
	id, err := app.snippets.Insert(form.Content, form.Title, form.Expires)

	if err != nil {
		app.serverError(w, err)
		return
	}
	//if the snippet was created successfully, we want to display a banner temporarily with text
	//"snippet created successfully" using session data
	//use the Put() method to add a string value ("Snippet successfully
	// created!") and the corresponding key ("flash") to the session data
	app.sessionManager.Put(r.Context(), "flash", "Snippet created successfully")
	//redirect to the created snippet after it has been created
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	//id, err := strconv.Atoi(r.URL.Query().Get("id"))

	//fetching parameters. returns a slice
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))

	if id < 1 || err != nil {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	fmt.Println(snippet)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	//wrapping our data in a struct
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.gohtml", data)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignUpForm{}
	app.render(w, http.StatusOK, "signup.gohtml", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {

	var form userSignUpForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "name cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "email cannot be blank")
	form.CheckField(validator.NotBlank(form.Name), "password", "password cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "password should be atleast 8 characters")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "enter a valid email address")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.gohtml", data)
		return
	}
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "email is already taken")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.gohtml", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked.
	app.sessionManager.Put(r.Context(), "flash", "Signed up successfully")
	//redirect user to login page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.gohtml", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Password), "password", "password cannot be empty")
	form.CheckField(validator.NotBlank(form.Email), "email", "email cannot be empty")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "password should be atleast 8 characters")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "enter a valid email address")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.gohtml", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Invalid email or password")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.gohtml", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//if email and password match then we create a new sessionID. Its good practice to generate a new
	//sessionID when authenticate state of privilege level changes for a user. Eg: after logging in or
	//logging out
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	//add the ID of the current user to session, so they are now "logged in"
	app.sessionManager.Put(r.Context(), "authenticatedUserId", id)

	//redirect user to create snippet page
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {

	//generate new token for user before logging the user out
	err := app.sessionManager.RenewToken(r.Context())

	if err != nil {
		app.serverError(w, err)
		return
	}

	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.sessionManager.Remove(r.Context(), "authenticatedUserId")

	//add a flash message indicating the user, that the user has been logged out
	app.sessionManager.Put(r.Context(), "flash", "Logged out successfully")

	//redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")
	user, err := app.users.Get(userId)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.User = user
	app.render(w, http.StatusOK, "account.gohtml", data)
}
