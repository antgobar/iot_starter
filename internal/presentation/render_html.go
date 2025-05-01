package presentation

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const baseTemplateDir = "templates"

var pageTemplates []string = []string{
	"home",
	"login",
	"register",
}

var baseTemplates []string = []string{
	"base",
	"navbar",
}

type CompliedTemplates map[string]*template.Template

type Templates struct {
	templates *CompliedTemplates
}

func NewHtmlPresenter() *Templates {
	templates, err := compileTemplates(baseTemplates, pageTemplates)
	if err != nil {
		log.Fatalln("Failed to compile templates:", err.Error())
	}
	return &Templates{templates}
}

func (t *Templates) Present(w http.ResponseWriter, r *http.Request, name string, payload any) error {
	tmpl, exists := (*t.templates)[name]
	if !exists {
		return errors.New(name + " template not found")
	}

	tmplName := "base"

	if r.Header.Get("HX-Request") == "true" {
		tmplName = "content"
	}

	return tmpl.ExecuteTemplate(w, tmplName, payload)
}

func compileTemplates(bases []string, pages []string) (*CompliedTemplates, error) {
	var templates = make(CompliedTemplates)
	for _, p := range pages {
		var allPages = make([]string, 0)
		for _, b := range bases {
			allPages = append(allPages, fmt.Sprintf("%s/%s.html", baseTemplateDir, b))
		}
		var err error
		allPages = append(allPages, fmt.Sprintf("%s/%s.html", baseTemplateDir, p))
		templates[p], err = template.ParseFiles(
			allPages...,
		)
		if err != nil {
			return nil, err
		}
	}
	return &templates, nil
}
