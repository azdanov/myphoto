package controllers

import (
	"fmt"
	"myphoto/views"
	"net/http"
)

// NewUsers creates a new Users Controller.
// This function will panic if templates are not correct,
// and should be used only during initial mux setup.
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("index", "users/new"),
	}
}

type Users struct {
	NewView *views.View
}

// New is used to render a form where a user can create an account.
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	err := u.NewView.Render(w, nil)
	if err != nil {
		panic(err)
	}
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create is used to for processing the user create account form.
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	fmt.Fprintln(w, form)
}
