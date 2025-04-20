package main

import (
	"context"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/measurement"
	"iotstarter/internal/store"
	"log"
	"time"
)

func main() {
	cfg, err := config.LoadTransformerConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	store, err := store.NewStore(ctx, cfg.DatabaseUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	handler := NewHandler(store)

	brokerClient, err := broker.NewBrokerClient(cfg.BrokerUrl)
	if err != nil {
		log.Println("ERROR: error connecting to broker client")
		return
	}
	defer brokerClient.Close()

	err = brokerClient.Subscribe(config.BrokerMeasurementSubject, handler.saveMeasurement)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("Transformer listening on subject: %s", config.BrokerMeasurementSubject)
	select {}
}

func NewHandler(store *store.Store) Handler {
	return Handler{store: store}
}

type Handler struct {
	store *store.Store
}

func (h Handler) saveMeasurement(m *measurement.Measurement) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()
	err := h.store.SaveMeasurement(ctx, *m)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Stored measurement under id", m.ID, "for device id", m.DeviceId)
}
