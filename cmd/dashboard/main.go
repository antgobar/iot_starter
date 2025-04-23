package main

import (
	"context"
	"iotstarter/internal/api"
	"iotstarter/internal/config"
	"iotstarter/internal/store"
	"iotstarter/internal/views"
	"log"
	"time"
)

func main() {
	dbUrl := config.MustLoadEnv("DATABASE_URL")
	apiAddr := config.MustLoadEnv("API_ADDR")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	store, err := store.NewPostgresStore(ctx, dbUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	views := views.NewViews()

	handler := api.NewHandler().WithStore(store).WithViews(views)
	server := api.NewServer(apiAddr, handler)
	server.Run("Dashboard")
}
