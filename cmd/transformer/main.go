package main

import (
	"context"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/consumer"
	"iotstarter/internal/store"
	"time"
)

func main() {
	dbUrl := config.MustLoadEnv("DATABASE_URL")
	brokerUrl := config.MustLoadEnv("BROKER_URL")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	store := store.NewPostgresStore(ctx, dbUrl)
	defer store.Close()

	brokerClient := broker.NewNatsBrokerClient(brokerUrl)
	defer brokerClient.Close()

	handler := consumer.NewHandler(store, brokerClient)
	handler.Run()
}
