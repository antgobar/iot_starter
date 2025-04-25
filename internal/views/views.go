package views

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const baseTemplateDir = "templates"

var htmlPages []string = []string{
	"home", "devices",
}

type CompliedTemplates map[string]*template.Template

type Views struct {
	templates *CompliedTemplates
}

type View struct {
	tmpl *template.Template
}

func (c *CompliedTemplates) Get(name string) *template.Template {
	tmpl, exists := (*c)[name]
	if !exists {
		log.Printf("Template %s not found", name)
		return nil
	}
	return tmpl
}

func (v *View) RenderTemplate(w http.ResponseWriter, r *http.Request, data any) error {
	var tmplName string

	if r.Header.Get("HX-Request") == "true" {
		tmplName = "content"
	} else {
		tmplName = "base"
	}

	return v.tmpl.ExecuteTemplate(w, tmplName, data)
}

func newView(tmpl *template.Template) View {
	return View{tmpl}
}

func NewViews() *Views {
	templates, err := compileTemplates(htmlPages)
	if err != nil {
		log.Fatalln("Failed to compile templates:", err.Error())
	}
	return &Views{templates}
}

func (v *Views) Page(name string) (View, error) {
	tmpl, exists := (*v.templates)[name]
	if !exists {
		return View{}, errors.New(name + " template not found")
	}
	return newView(tmpl), nil
}

func compileTemplates(pages []string) (*CompliedTemplates, error) {
	var templates = make(CompliedTemplates)
	for _, p := range pages {
		var err error
		templates[p], err = template.ParseFiles(
			baseTemplateDir+"/base.html",
			fmt.Sprintf("%s/%s.html", baseTemplateDir, p),
		)
		if err != nil {
			return nil, err
		}
	}
	return &templates, nil
}
