package main

import (
	"context"
	"iotstarter/internal/api"
	"iotstarter/internal/auth"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/consumer"
	"iotstarter/internal/database"
	"iotstarter/internal/device"
	"iotstarter/internal/gateway"
	"iotstarter/internal/measurement"
	"iotstarter/internal/middleware"
	"iotstarter/internal/presentation"
	"iotstarter/internal/session"
	"iotstarter/internal/user"
	"iotstarter/internal/web"
	"time"
)

func main() {
	dbUrl := config.MustLoadEnv("DATABASE_URL")
	apiAddr := config.MustLoadEnv("API_ADDR")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	brokerClient := broker.NewMemoryBroker()

	db := database.NewPostgresPool(ctx, dbUrl)

	userRepo := user.NewPostgresRepository(ctx, db.Pool)
	sessionRepo := session.NewPostgresRepository(ctx, db.Pool)
	deviceRepo := device.NewPostgresRepository(ctx, db.Pool)
	measurementRepo := measurement.NewPostgresRepository(ctx, db.Pool)

	htmlPresenter := presentation.NewHtmlPresenter()
	userService := user.NewService(userRepo)
	sessionService := session.NewService(sessionRepo)
	authService := auth.NewService(userRepo, sessionRepo)
	deviceService := device.NewService(deviceRepo)
	measurementService := measurement.NewService(measurementRepo)
	devicesRepo := device.NewPostgresRepository(ctx, db.Pool)

	devicesService := device.NewService(devicesRepo)
	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(authService)
	deviceHandler := device.NewHandler(deviceService, htmlPresenter)
	webPageHandler := web.NewHandler(htmlPresenter)
	measurementHandler := measurement.NewHandler(measurementService)

	gatewayService := gateway.NewService(brokerClient)
	gatewayHandler := gateway.NewHandler(gatewayService, devicesService, config.BrokerMeasurementSubject)
	consumerService := consumer.NewService(brokerClient, measurementRepo, config.BrokerMeasurementSubject)

	consumerHandler := consumer.NewHandler(consumerService)
	consumerHandler.Run()

	middlewareStack := middleware.LoadMiddleware(sessionService)
	server := api.NewServer(
		apiAddr,
		middlewareStack,
		authHandler,
		userHandler,
		deviceHandler,
		webPageHandler,
		measurementHandler,
		gatewayHandler,
	)
	server.Run("Monolith")

	select {}
}
