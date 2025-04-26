package api

import (
	"iotstarter/internal/middleware"
	"log"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(addr string, h *Handler) Server {
	stack := middleware.LoadMiddleware(h.store)
	mux := h.registerUserRoutes()
	server := &http.Server{
		Addr:    addr,
		Handler: stack(mux),
	}
	return Server{server: server}
}

func (s Server) Run(appName string) {
	log.Println(appName, "starting on", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v", s.server.Addr, err)
	}
}
