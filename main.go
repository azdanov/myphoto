package main

import (
	"fmt"
	"myphoto/controllers"
	"myphoto/middleware"
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
	_ = svc.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(svc.User)
	galleriesC := controllers.NewGalleries(svc.Gallery)

	requireUserMw := middleware.RequireUser{UserService: svc.User}

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.LoginUser).Methods("POST")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.CreateUser).Methods("POST")

	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")

	fmt.Println("Starting on: http://localhost:3000")
	if err := http.ListenAndServe("localhost:3000", r); err != nil {
		panic(err)
	}
}
