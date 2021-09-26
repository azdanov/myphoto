package controllers

import (
	"myphoto/models"
	"myphoto/views"
	"net/http"
)

// NewUsers creates a new Users Controller.
// This function will panic if templates are not correct,
// and should be used only during initial mux setup.
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("index", "users/new"),
		us:      us,
	}
}

type Users struct {
	NewView *views.View
	us      *models.UserService
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
	Name     string `schema:"name"`
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

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password, // Unsafe
	}

	err := u.us.Create(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
	}
}
