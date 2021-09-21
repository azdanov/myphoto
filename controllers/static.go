package controllers

import "myphoto/views"

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("index", "static/home"),
		Contact: views.NewView("index", "static/contact"),
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
}
