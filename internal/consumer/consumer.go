package consumer

import (
	"context"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/model"
	"iotstarter/internal/store"
	"log"
	"time"
)

type Handler struct {
	store     *store.Store
	broker    broker.Broker
	consumers []string
}

func NewHandler(store *store.Store, broker broker.Broker) *Handler {
	return &Handler{store, broker, nil}
}

func (h *Handler) registerConsumers() {
	err := h.broker.Subscribe(config.BrokerMeasurementSubject, h.saveMeasurement)
	if err != nil {
		log.Fatalln(err.Error())
	}
	h.consumers = append(h.consumers, config.BrokerMeasurementSubject)
}

func (h *Handler) saveMeasurement(m *model.Measurement) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()
	err := h.store.SaveMeasurement(ctx, m)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Stored measurement under id", m.ID, "for device id", m.DeviceId)
}

func (h *Handler) Run() {
	h.registerConsumers()
	log.Printf("Transformer listening on subject(s): %s", h.consumers)
	select {}
}
