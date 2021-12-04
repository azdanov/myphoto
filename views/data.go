package views

import (
	"errors"
	"html/template"
	"log"
	"myphoto/models"
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
