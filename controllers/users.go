package controllers

import (
	"fmt"
	"myphoto/models"
	"myphoto/views"
	"net/http"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

// NewUsers creates a new Users Controller.
// This function will panic if templates are not correct,
// and should be used only during initial mux setup.
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("index", "users/new"),
		LoginView: views.NewView("index", "users/login"),
		us:        us,
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

	user := toUserModel(form)
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
	}
}

func toUserModel(form SignupForm) models.User {
	return models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login is used to authenticate a user.
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	fmt.Fprintln(w, form)
}
