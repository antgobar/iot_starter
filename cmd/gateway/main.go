package main

import (
	"iotstarter/internal/api"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
)

func main() {
	brokerUrl := config.MustLoadEnv("BROKER_URL")
	gatewayAddr := config.MustLoadEnv("GATEWAY_ADDR")

	brokerClient := broker.NewNatsBrokerClient(brokerUrl)
	defer brokerClient.Close()

	handler := api.NewHandler().WithBroker(brokerClient)
	server := api.NewServer(gatewayAddr, handler)
	server.Run("Gateway")
}
