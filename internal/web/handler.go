package web

import (
	"iotstarter/internal/presentation"
	"log"
	"net/http"
)

type Handler struct {
	presenter presentation.Presenter
}

func NewHandler(p presentation.Presenter) *Handler {
	return &Handler{presenter: p}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.home)
	mux.HandleFunc("/register", h.register)
	mux.HandleFunc("/login", h.login)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.HandlerFunc(staticResources)))
	mux.HandleFunc("/favicon.ico", favicon)
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := h.presenter.Present(w, r, "home", nil); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if err := h.presenter.Present(w, r, "login", nil); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if err := h.presenter.Present(w, r, "register", nil); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func staticResources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	fs := http.FileServer(http.Dir("static"))
	fs.ServeHTTP(w, r)
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/favicon.ico", http.StatusMovedPermanently)
}
