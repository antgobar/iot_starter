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
	dbUrl := config.MustLoadEnv("DATABASE_URL")
	brokerUrl := config.MustLoadEnv("BROKER_URL")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()
	store, err := store.NewStore(ctx, dbUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	brokerClient, err := broker.NewNatsBrokerClient(brokerUrl)
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
