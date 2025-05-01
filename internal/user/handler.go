package user

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", h.register)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	username := r.FormValue("username")
	password := r.FormValue("password")
	err := h.svc.Register(ctx, username, password)

	if err == ErrUsernameTaken {
		http.Error(w, "username taken", http.StatusConflict)
		return
	}

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "error registering user", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
