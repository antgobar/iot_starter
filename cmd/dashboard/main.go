package main

import (
	"context"
	"iotstarter/internal/api"
	"iotstarter/internal/config"
	"iotstarter/internal/device"
	"iotstarter/internal/middleware"
	"iotstarter/internal/store"
	"iotstarter/internal/views"
	"time"
)

func main() {
	dbUrl := config.MustLoadEnv("DATABASE_URL")
	apiAddr := config.MustLoadEnv("API_ADDR")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	store := store.NewPostgresStore(ctx, dbUrl)
	defer store.Close()

	deviceRepo := device.NewPostgresRepository(ctx, dbUrl)
	deviceService := device.NewService(deviceRepo)
	deviceHandler := device.NewHandler(deviceService)

	views := views.NewViews()

	_ = api.NewHandler().WithViewsAndStore(views, store)

	middlewareStack := middleware.LoadMiddleware(store)
	server := api.NewServer(apiAddr, middlewareStack, deviceHandler)
	server.Run("Dashboard")
}
