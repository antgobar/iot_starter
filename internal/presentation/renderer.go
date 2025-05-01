package presentation

import "net/http"

type Presenter interface {
	Present(w http.ResponseWriter, r *http.Request, name string, data any) error
}
