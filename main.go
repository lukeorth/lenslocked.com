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
    must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    must(contactView.Render(w, nil))
}

func main() {
    homeView = views.NewView("bootstrap", "views/home.html")
    contactView = views.NewView("bootstrap", "views/contact.html")

    r := mux.NewRouter()
    r.HandleFunc("/", home)
    r.HandleFunc("/contact", contact)
    must(http.ListenAndServe(":3000", r))
}

func must(err error) {
    if err != nil {
        panic(err)
    }
}
