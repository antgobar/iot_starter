package presentation

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	baseTemplateDir = "templates"
	basesCount      = 2
	pagesCount      = 5
	fragmentsCount  = 1
)

var bases [basesCount]string = [basesCount]string{
	"base",
	"navbar",
}

var pages [pagesCount]string = [pagesCount]string{
	"home",
	"login",
	"register",
	"devices",
	"device",
}

var fragments [fragmentsCount]string = [fragmentsCount]string{
	"device-reauth",
}

type CompliedTemplates map[string]*template.Template

type Templates struct {
	templates *CompliedTemplates
}

func NewHtmlPresenter() *Templates {
	templates, err := compileTemplates(bases, pages, fragments)
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

func compileTemplates(
	bases [basesCount]string,
	pages [pagesCount]string,
	fragments [fragmentsCount]string,
) (*CompliedTemplates, error) {

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

	var err error
	for _, f := range fragments {
		templates[f], err = template.ParseFiles(fmt.Sprintf("%s/%s.html", baseTemplateDir, f))
		if err != nil {
			return nil, err
		}
	}

	return &templates, nil
}
