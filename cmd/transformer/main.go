package main

import (
	"iotstarter/internal/broker"
	"iotstarter/internal/measurement"
	"log"
	"os"
)

func main() {
	natsSubject := os.Getenv("BROKER_SUBJECT")
	brokerUrl := os.Getenv("BROKER_URL")

	log.Printf("Worker listening on subject: %s", natsSubject)

	brokerClient, err := broker.NewBrokerClient(brokerUrl)
	if err != nil {
		log.Println("ERROR: error connecting to broker client")
		return
	}
	defer brokerClient.Close()

	brokerClient.Subscribe(natsSubject, func(measurement *measurement.Measurement) {
		log.Println("Received message:", *measurement)
	})

	select {}
}
