package main

import (
	"context"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/store"
	"iotstarter/internal/transformer"
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

	brokerClient, err := broker.NewNatsBrokerClient(cfg.BrokerUrl)
	if err != nil {
		log.Println("ERROR: error connecting to broker client")
		return
	}
	defer brokerClient.Close()

	handler := transformer.NewHandler(store)

	err = brokerClient.Subscribe(config.BrokerMeasurementSubject, handler.SaveMeasurement)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("Transformer listening on subject: %s", config.BrokerMeasurementSubject)
	select {}
}
