package controllers

import (
	"fmt"
	"myphoto/models"
	"myphoto/views"
	"net/http"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("index", "galleries/new"),
		// IndexView: views.NewView("index", "galleries/index"),
		// ShowView:  views.NewView("index", "galleries/show"),
		// EditView:  views.NewView("index", "galleries/edit"),
		gs: gs,
	}
}

type Galleries struct {
	New       *views.View
	IndexView *views.View
	ShowView  *views.View
	EditView  *views.View
	gs        models.GalleryService
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// Create is used to for processing the gallery create account form.
// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}
	gallery := models.Gallery{
		Title: form.Title,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}
