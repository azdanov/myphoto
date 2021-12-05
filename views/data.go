package views

import (
	"errors"
	"html/template"
	"log"
	"myphoto/models"
	"net/http"
	"time"
)

const (
	// AlertLevelError level corresponds to bootstrap danger class.
	AlertLevelError = "danger"
	// AlertLevelWarning level corresponds to bootstrap warning class.
	AlertLevelWarning = "warning"
	// AlertLevelInfo level corresponds to bootstrap info class.
	AlertLevelInfo = "info"
	// AlertLevelSuccess level corresponds to bootstrap success class.
	AlertLevelSuccess = "success"
	// AlertMessageGeneric is displayed when any random error is encountered.
	AlertMessageGeneric = "Something went wrong. Please try again, and contact us if the problem persists."
)

// Alert is used to render notifications in views.
type Alert struct {
	Level   string
	Message string
}

// Data is the top level structure that views accept.
type Data struct {
	Alert *Alert
	User  *models.User
	CSRF  template.HTML
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	var publicError PublicError
	if ok := errors.As(err, &publicError); ok {
		d.Alert = &Alert{
			Level:   AlertLevelError,
			Message: publicError.Public(),
		}
	} else {
		log.Println(err)
		d.Alert = &Alert{
			Level:   AlertLevelError,
			Message: AlertMessageGeneric,
		}
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLevelError,
		Message: msg,
	}
}

type PublicError interface {
	error
	Public() string
}

func persistAlert(w http.ResponseWriter, alert Alert) {
	expiresAt := time.Now().Add(5 * time.Minute)
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    alert.Level,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    alert.Message,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

func clearAlert(w http.ResponseWriter) {
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

func getAlert(r *http.Request) *Alert {
	lvl, err := r.Cookie("alert_level")
	if err != nil {
		return nil
	}
	msg, err := r.Cookie("alert_message")
	if err != nil {
		return nil
	}
	alert := Alert{
		Level:   lvl.Value,
		Message: msg.Value,
	}
	return &alert
}

// RedirectAlert extends http.Redirect and persists an alert in a cookie, so it can be displayed when a new page is loaded.
func RedirectAlert(w http.ResponseWriter, r *http.Request, url string, code int, alert Alert) {
	persistAlert(w, alert)
	http.Redirect(w, r, url, code)
}
