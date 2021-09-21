package controllers

import "myphoto/views"

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("index", "views/static/home.gohtml"),
		Contact: views.NewView("index", "views/static/contact.gohtml"),
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
}
