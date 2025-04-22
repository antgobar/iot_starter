package views

import (
	"net/http"
	"text/template"
)

const baseTemplateDir = "templates/"

type Views struct{}

type View struct {
	w    http.ResponseWriter
	tmpl *template.Template
}

func (v *View) Render(data any) {
	v.tmpl.Execute(v.w, data)
}

func newView(w http.ResponseWriter, tmpl *template.Template) View {
	return View{w, tmpl}
}

func NewViews() *Views { return &Views{} }

func (v *Views) IndexPage(w http.ResponseWriter) View {
	tmpl := compileTemplate(baseTemplateDir + "index.html")
	return newView(w, tmpl)
}

func compileTemplate(path string) *template.Template {
	return template.Must(template.ParseFiles(path))
}
