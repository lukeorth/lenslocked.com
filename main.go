package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/lukeorth/lenslocked.com/controllers"
	"github.com/lukeorth/lenslocked.com/middleware"
	"github.com/lukeorth/lenslocked.com/models"
	"github.com/lukeorth/lenslocked.com/rand"
)

const (
    host = "localhost"
    port = 5432
    user = "postgres"
    dbname = "lenslocked_dev"
)

func main() {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
    services, err := models.NewServices(psqlInfo)
    if err != nil {
        panic(err)
    }
    must(err)
    services.AutoMigrate()

    r := mux.NewRouter()
    staticC := controllers.NewStatic()
    usersC := controllers.NewUsers(services.User)
    galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

    // TODO: Update this to be a config variable
    isProd := false
    b, err := rand.Bytes(32)
    must(err)
    csrfMw := csrf.Protect(b, csrf.Secure(isProd))
    userMw := middleware.User{
        UserService: services.User,
    }
    requireUserMw := middleware.RequireUser{
        User: userMw,
    }

    r.Handle("/", staticC.Home).Methods("GET")
    r.Handle("/contact", staticC.Contact).Methods("GET")
    r.HandleFunc("/signup", usersC.New).Methods("GET")
    r.Handle("/signup", usersC.NewView).Methods("GET")
    r.HandleFunc("/signup", usersC.Create).Methods("POST")
    r.Handle("/login", usersC.LoginView).Methods("GET")
    r.HandleFunc("/login", usersC.Login).Methods("POST")

    // Assets
    assetHandler := http.FileServer(http.Dir("./assets/"))
    assetHandler = http.StripPrefix("/assets/", assetHandler)
    r.PathPrefix("/assets/").Handler(assetHandler)

    // Image routes
    imageHandler := http.FileServer(http.Dir("./images/"))
    r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

    // Gallery routes
    r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
    r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
    r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
    r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
    r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
    r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
    r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")

    // /galleries/:id/images/:filename/delete
    r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")
    r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)

    must(http.ListenAndServe(":3000", csrfMw(userMw.Apply(r))))
}

func must(err error) {
    if err != nil {
        panic(err)
    }
}
