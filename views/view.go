package views

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

var (
	layoutDir   = "views/layouts/"
	templateDir = "views/"
	templateExt = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	switch data.(type) {
	case Data:
		// empty
	default:
		data = Data{
			Yield: data,
		}
	}
	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, data); err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. If the problem persists, please email support@myphoto.com.", http.StatusInternalServerError)
		return
	}
	_, err := io.Copy(w, &buf)
	if err != nil {
		panic(err)
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(layoutDir + "*" + templateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath takes a slice of string
// representing file paths for templates
// to prepend templateDir directory to each.
//
// The argument files with {"home"} would result
// in {"views/home"} if templateDir is "views/".
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = templateDir + f
	}
}

// addTemplateExt takes a slice of string
// representing file paths for templates
// to append templateExt file extension to each.
//
// The argument files with {"home"} would result
// in {"home.gohtml"} if templateExt is ".gohtml".
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + templateExt
	}
}
