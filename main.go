package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lukeorth/lenslocked.com/views"
)

var (
    homeView *views.View
    contactView *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    err := homeView.Template.Execute(w, nil)
    if err != nil {
        panic(err)
    }
}

func contact(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    err := contactView.Template.Execute(w, nil)
    if err != nil {
        panic(err)
    }
}

func main() {
    homeView = views.NewView("views/home.html")
    contactView = views.NewView("views/contact.html")
    r := mux.NewRouter()
    r.HandleFunc("/", home)
    r.HandleFunc("/contact", contact)
    http.ListenAndServe(":3000", r)
}
