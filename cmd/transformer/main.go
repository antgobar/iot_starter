package main

import (
	"context"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/consumer"
	"iotstarter/internal/measurement"
	"iotstarter/internal/postgres"
	"time"
)

func main() {
	dbUrl := config.MustLoadEnv("DATABASE_URL")
	brokerUrl := config.MustLoadEnv("BROKER_URL")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	brokerClient := broker.NewNatsBrokerClient(brokerUrl)
	defer brokerClient.Close()

	db := postgres.NewPostgresPool(ctx, dbUrl)

	measurementRepo := measurement.NewPostgresRepository(ctx, db.Pool)
	svc := consumer.NewService(brokerClient, measurementRepo, config.BrokerMeasurementSubject)

	handler := consumer.NewHandler(svc)
	handler.Run()
	select {}
}
