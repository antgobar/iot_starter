package main

import (
	"context"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/consumer"
	"iotstarter/internal/store"
	"log"
	"time"
)

func main() {
	dbUrl := config.MustLoadEnv("DATABASE_URL")
	brokerUrl := config.MustLoadEnv("BROKER_URL")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()
	store, err := store.NewPostgresStore(ctx, dbUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	brokerClient, err := broker.NewNatsBrokerClient(brokerUrl)
	if err != nil {
		log.Println("ERROR: error connecting to broker client")
		return
	}
	defer brokerClient.Close()

	handler := consumer.NewHandler(store, brokerClient)
	handler.Run()
}
