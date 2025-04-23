package main

import (
	"iotstarter/internal/api"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/logging"
	"log"
)

func main() {
	logging.SetUp()
	brokerUrl := config.MustLoadEnv("BROKER_URL")
	gatewayAddr := config.MustLoadEnv("GATEWAY_ADDR")

	brokerClient, err := broker.NewNatsBrokerClient(brokerUrl)
	if err != nil {
		log.Fatalln("ERROR: ", err.Error())
	}
	defer brokerClient.Close()

	handler := api.NewHandler().WithBroker(brokerClient)
	server := api.NewServer(gatewayAddr, handler)
	server.Run("Gateway")
}
