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
	broker *broker.Broker
}

func newHandler(store *store.Store, broker *broker.Broker) Handler {
	return Handler{store: store, broker: broker}
}

func registerUserRoutes(h Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /devices", h.registerDevice)
	mux.HandleFunc("GET /devices", h.getDevices)
	mux.HandleFunc("GET /devices/{id}/measurements", h.getDeviceMeasurements)
	// mux.HandleFunc("POST /measurements", h.saveMeasurement)
	return mux
}

type Server struct {
	server *http.Server
}

func NewServer(addr string, store *store.Store, broker *broker.Broker) Server {
	stack := middleware.LoadMiddleware()
	handler := newHandler(store, broker)
	mux := registerUserRoutes(handler)
	server := &http.Server{
		Addr:    addr,
		Handler: stack(mux),
	}
	return Server{server: server}
}

func (s Server) WithBroker(broker broker.BrokerClient) {

}

func (s Server) Run(appName string) {
	log.Println(appName+"starting on", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v", s.server.Addr, err)
	}
}
