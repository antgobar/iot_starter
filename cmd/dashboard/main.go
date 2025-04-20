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
	cfg, err := config.LoadDashboardConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}
	if err != nil {
		log.Fatalln(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	store, err := store.NewStore(ctx, cfg.DatabaseUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	handler := api.NewHandler(store, nil)
	server := api.NewServer(cfg.Addr, handler)
	server.Run("Dashboard")
}
