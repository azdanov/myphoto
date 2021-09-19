package main

import (
	"myphoto/views"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	homeTemplate    *views.View
	contactTemplate *views.View
)

func home(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := homeTemplate.Template.ExecuteTemplate(w, homeTemplate.Layout, nil)
	if err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := contactTemplate.Template.ExecuteTemplate(w, contactTemplate.Layout, nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	homeTemplate = views.NewView("index", "views/home.gohtml")
	contactTemplate = views.NewView("index", "views/contact.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	http.ListenAndServe("localhost:3000", r)
}
