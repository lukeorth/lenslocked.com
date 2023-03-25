package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lukeorth/lenslocked.com/controllers"
	"github.com/lukeorth/lenslocked.com/models"
)

const (
    host = "localhost"
    port = 5432
    user = "postgres"
    dbname = "lenslocked_dev"
)

func main() {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
    us, err := models.NewUserService(psqlInfo)
    if err != nil {
        panic(err)
    }
    must(err)
    us.AutoMigrate()

    staticC := controllers.NewStatic()
    usersC := controllers.NewUsers(us)

    r := mux.NewRouter()
    r.Handle("/", staticC.Home).Methods("GET")
    r.Handle("/contact", staticC.Contact).Methods("GET")
    r.HandleFunc("/signup", usersC.New).Methods("GET")
    r.Handle("/signup", usersC.NewView).Methods("GET")
    r.HandleFunc("/signup", usersC.Create).Methods("POST")
    r.Handle("/login", usersC.LoginView).Methods("GET")
    r.HandleFunc("/login", usersC.Login).Methods("POST")
    r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
    must(http.ListenAndServe(":3000", r))
}

func must(err error) {
    if err != nil {
        panic(err)
    }
}
