package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
	"github.com/lukeorth/lenslocked.com/context"
)

var (
    LayoutDir string = "views/layouts/"
    TemplateDir string = "views/"
    TemplateExt string = ".html"
)

func NewView(layout string, files ...string) *View {
    addTemplatePath(files)
    addTemplateExt(files)
    files = append(files, layoutFiles()...)
    t, err := template.New("").Funcs(template.FuncMap{
        "csrfField": func() (template.HTML, error) {
            return "", errors.New("csrfField is not implemented")
        },
    }).ParseFiles(files...)
    if err != nil {
        panic(err)
    }

    return &View{
        Template: t,
        Layout: layout,
    }
}

type View struct {
    Template *template.Template
    Layout string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    v.Render(w, r, nil)
}

// Render is used to render the view with the predefined layout
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
    w.Header().Set("Content-Type", "text/html")
    var vd Data
    switch d := data.(type) {
    case Data:
        vd = d
        // do nothing
    default:
        vd = Data{
            Yield: data,
        }
    }
    vd.User = context.User(r.Context())
    var buf bytes.Buffer
    csrfField := csrf.TemplateField(r)
    tpl := v.Template.Funcs(template.FuncMap{
        "csrfField": func() template.HTML {
            return csrfField
        },
    })
    if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
        log.Println(err)
        http.Error(w, "Something went wrong, If the problem persists, please email support@lenslocked.com", http.StatusInternalServerError)
        return
    }
    io.Copy(w, &buf)
}

// layoutFiles returns a slice of strings representing
// the layout files used in our application.
func layoutFiles() []string {
    files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
    if err != nil {
        panic(err)
    }
    return files
}

// addtTemplatePath takes in a slice of strings
// representing file paths for templates, and it 
// prepends the TemplateDir directory to each string
// in the slice 
// 
// Eg the input {"home"} would result in the output
// {"views/home"} if TemplateDir == "views/"
func addTemplatePath(files []string) {
    for i, f := range files {
        files[i] = TemplateDir + f
    }
}

// addTemplateExt takes in a slice of strings 
// representing file paths for templates and it 
// appends the TemplateExt extension to each 
// string in the slice 
// 
// Eg the input {"home"} would result in the output
// {"home.html"} if the TemplateExt == ".html"
func addTemplateExt(files []string) {
    for i, f := range files {
        files[i] = f + TemplateExt
    }
}
