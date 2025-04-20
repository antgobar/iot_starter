package main

import (
	"iotstarter/internal/api"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadGatewayConfig()
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
	}
	brokerClient, err := broker.NewNatsBrokerClient(cfg.BrokerUrl)
	if err != nil {
		log.Fatalln("ERROR: ", err.Error())
	}
	defer brokerClient.Close()

	handler := api.NewHandler(nil, brokerClient)
	server := api.NewServer(cfg.Addr, handler)
	server.Run("Gateway")
}
