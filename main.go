package main

import (
	"fmt"
	"myphoto/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	fmt.Println("Starting on: http://localhost:3000")
	if err := http.ListenAndServe("localhost:3000", r); err != nil {
		return
	}
}
