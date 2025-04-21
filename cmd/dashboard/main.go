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
	dbUrl := config.MustLoadEnv("DATABASE_URL")
	apiAddr := config.MustLoadEnv("API_ADDR")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	store, err := store.NewStore(ctx, dbUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	handler := api.NewHandler().WithStore(store)
	server := api.NewServer(apiAddr, handler)
	server.Run("DashboardApi")
}
