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
	apiAddr := config.MustLoadEnv("API_ADDR")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	store, err := store.NewStore(ctx, dbUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	brokerClient := broker.NewMemoryBroker()
	apiHandler := api.NewHandler().WithStore(store).WithBroker(brokerClient)

	server := api.NewServer(apiAddr, apiHandler)
	go server.Run("IOT Monolith")

	consumerHandler := consumer.NewHandler(store, brokerClient)
	consumerHandler.Run()

	select {}
}
