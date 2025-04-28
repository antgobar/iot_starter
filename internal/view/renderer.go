package view

import "net/http"

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, name string, data any) error
}
