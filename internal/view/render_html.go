package view

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const baseTemplateDir = "templates"

var htmlPages []string = []string{
	"home",
	"devices",
	"login",
	"register",
}

type CompliedTemplates map[string]*template.Template

type Views struct {
	templates *CompliedTemplates
}

func NewHtmlView() *Views {
	templates, err := compileTemplates(htmlPages)
	if err != nil {
		log.Fatalln("Failed to compile templates:", err.Error())
	}
	return &Views{templates}
}

func (v *Views) Render(w http.ResponseWriter, r *http.Request, name string, data any) error {
	tmpl, exists := (*v.templates)[name]
	if !exists {
		return errors.New(name + " template not found")
	}

	var tmplName string

	if r.Header.Get("HX-Request") == "true" {
		tmplName = "content"
	} else {
		tmplName = "base"
	}

	return tmpl.ExecuteTemplate(w, tmplName, data)
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
