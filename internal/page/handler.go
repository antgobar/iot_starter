package page

import (
	"iotstarter/internal/view"
	"log"
	"net/http"
)

type Handler struct {
	renderer view.Renderer
}

func NewHandler(v view.Renderer) *Handler {
	return &Handler{renderer: v}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.home)
	mux.HandleFunc("/register", h.register)
	mux.HandleFunc("/login", h.login)
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := h.renderer.Render(w, r, "home", nil); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if err := h.renderer.Render(w, r, "login", nil); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		log.Println("ERROR:", err.Error())
	}
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if err := h.renderer.Render(w, r, "register", nil); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		log.Println("ERROR:", err.Error())
	}
}
