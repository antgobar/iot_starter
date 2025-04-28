package presentation

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct{}

func NewJsonPresentation() *JsonResponse {
	return &JsonResponse{}
}

func (j *JsonResponse) Present(w http.ResponseWriter, r *http.Request, name string, data any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}
