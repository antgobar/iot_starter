package api

import (
	"iotstarter/internal/broker"
	"iotstarter/internal/middleware"
	"iotstarter/internal/store"
	"log"
	"net/http"
)

type Handler struct {
	store  *store.Store
	broker broker.Broker
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) WithBroker(broker broker.Broker) *Handler {
	h.broker = broker
	return h
}

func (h *Handler) WithStore(store *store.Store) *Handler {
	h.store = store
	return h
}

type Server struct {
	server *http.Server
}

func NewServer(addr string, handler *Handler) Server {
	stack := middleware.LoadMiddleware()
	mux := registerUserRoutes(handler)
	server := &http.Server{
		Addr:    addr,
		Handler: stack(mux),
	}
	return Server{server: server}
}

func (s Server) Run(appName string) {
	log.Println(appName+"starting on", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v", s.server.Addr, err)
	}
}
