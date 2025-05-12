package auth

import (
	"context"
	"iotstarter/internal/session"
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
	mux.HandleFunc("POST /login", h.logIn)
	mux.HandleFunc("GET /logout", h.logOut)
}

func (h *Handler) logIn(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	username := r.FormValue("username")
	password := r.FormValue("password")

	sesh, err := h.svc.LogIn(ctx, username, password)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "error logging in", http.StatusInternalServerError)
		return
	}

	session.SetCookie(w, string(sesh.Token))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := UserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := h.svc.LogOut(ctx, user.ID)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
	session.ClearCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
