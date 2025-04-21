package main

import (
	"context"
	"iotstarter/internal/api"
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
	apiAddr := config.MustLoadEnv("API_ADDR")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	store, err := store.NewStore(ctx, dbUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	brokerClient, err := broker.NewNatsBrokerClient(brokerUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	apiHandler := api.NewHandler(store, brokerClient)

	server := api.NewServer(apiAddr, apiHandler)
	go server.Run("IOT Monolith")

	consumerHandler := consumer.NewHandler(store, brokerClient)
	consumerHandler.Run()

	select {}
}
