package main

import (
	"context"
	"iotstarter/internal/api"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/database"
	"iotstarter/internal/device"
	"iotstarter/internal/gateway"
	"iotstarter/internal/middleware"
	"time"
)

func main() {
	brokerUrl := config.MustLoadEnv("BROKER_URL")
	gatewayAddr := config.MustLoadEnv("GATEWAY_ADDR")
	dbUrl := config.MustLoadEnv("DATABASE_URL")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	db := database.NewPostgresPool(ctx, dbUrl)

	brokerClient := broker.NewNatsBrokerClient(brokerUrl)
	defer brokerClient.Close()

	devicesRepo := device.NewPostgresRepository(ctx, db.Pool)
	devicesService := device.NewService(devicesRepo)

	gatewayService := gateway.NewService(brokerClient)
	gatewayHandler := gateway.NewHandler(gatewayService, devicesService, config.BrokerMeasurementSubject)

	middlewareStack := middleware.LoadLoggingMiddleware()
	server := api.NewServer(
		gatewayAddr,
		middlewareStack,
		gatewayHandler,
	)
	server.Run("Gateway")
}
