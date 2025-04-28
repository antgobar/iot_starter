package main

import (
	"iotstarter/internal/api"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/gateway"
	"iotstarter/internal/middleware"
)

func main() {
	brokerUrl := config.MustLoadEnv("BROKER_URL")
	gatewayAddr := config.MustLoadEnv("GATEWAY_ADDR")

	brokerClient := broker.NewNatsBrokerClient(brokerUrl)
	defer brokerClient.Close()

	gatewayService := gateway.NewService(brokerClient)
	gatewayHandler := gateway.NewHandler(gatewayService, config.BrokerMeasurementSubject)

	middlewareStack := middleware.LoadMiddleware(nil)
	server := api.NewServer(
		gatewayAddr,
		middlewareStack,
		gatewayHandler,
	)
	server.Run("Gateway")
}
