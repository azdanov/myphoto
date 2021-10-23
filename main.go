package main

import (
	"fmt"
	"myphoto/controllers"
	"myphoto/models"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "myphoto"
	sslmode  = "disable"
	timeZone = "Europe/Tallinn"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, port, user, password, dbname, sslmode, timeZone)
	svc, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer svc.Close()
	svc.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(svc.User)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.LoginUser).Methods("POST")
	r.HandleFunc("/signup", usersC.Create).Methods("GET")
	r.HandleFunc("/signup", usersC.CreateUser).Methods("POST")

	fmt.Println("Starting on: http://localhost:3000")
	if err := http.ListenAndServe("localhost:3000", r); err != nil {
		panic(err)
	}
}
