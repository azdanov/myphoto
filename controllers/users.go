package controllers

import (
	"errors"
	"myphoto/context"
	"myphoto/models"
	"myphoto/rand"
	"myphoto/views"
	"net/http"
	"time"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

// NewUsers creates a new Users Controller.
// This function will panic if templates are not correct,
// and should be used only during initial mux setup.
func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("index", "users/new"),
		LoginView: views.NewView("index", "users/login"),
		us:        us,
	}
}

// New is used to render the form where a user can
// create a new user account.
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create is used to for processing the user create account form.
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	user := toUserModel(form)
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	if err := u.signIn(w, &user); err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
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
	vd := views.Data{}
	var form LoginForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrResourceNotFound), errors.Is(err, models.ErrInvalidPassword):
			vd.AlertError("Invalid email address or password")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}

	if err = u.signIn(w, user); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// signIn is used to sign the user in by creating a cookie.
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

// Logout is used to delete a users' session cookie (remember_token)
// and then will update the user resource with a new remember_token.
// POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	user := context.User(r.Context())
	token, _ := rand.RememberToken()
	user.Remember = token
	_ = u.us.Update(user)
	http.Redirect(w, r, "/", http.StatusFound)
}
