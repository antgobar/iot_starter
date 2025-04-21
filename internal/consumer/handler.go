package consumer

import (
	"iotstarter/internal/broker"
	"iotstarter/internal/store"
	"log"
)

type Handler struct {
	store     *store.Store
	broker    broker.Broker
	consumers []Consumer
}

func NewHandler(store *store.Store, broker broker.Broker) *Handler {
	return &Handler{store, broker, nil}
}

func (h *Handler) consumersSubjects() []string {
	subjects := make([]string, 0)
	for _, consumer := range h.consumers {
		subjects = append(subjects, consumer.subject)
	}
	return subjects
}

func (h *Handler) Run() {
	h.registerConsumers()
	log.Printf("Transformer listening on subject(s): %s", h.consumersSubjects())
	select {}
}
