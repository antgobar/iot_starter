package main

import (
	"context"
	"iotstarter/internal/api"
	"iotstarter/internal/config"
	"iotstarter/internal/store"
	"log"
	"time"
)

func main() {
	config, err := config.LoadDashboardConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}
	if err != nil {
		log.Fatalln(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	store, err := store.NewStore(ctx, config.DatabaseUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	serverA := api.NewServer(config.Addr, store, nil)
	serverA.Run("Dashboard")
}
