package main

import (
	"encoding/json"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/measurement"
	"iotstarter/internal/middleware"
	"log"
	"net/http"
)

func main() {
	config, err := config.LoadGatewayConfig()
	if err != nil {
		log.Println("ERROR:", err.Error())
		return
	}
	brokerClient, err := broker.NewBrokerClient(config.BrokerUrl)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		return
	}
	defer brokerClient.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /measurement/schema", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		example := measurement.Measurement{
			DeviceId: 12345,
			Name:     "temperature",
			Value:    50,
			Unit:     "C",
		}
		json.NewEncoder(w).Encode(example)
	})

	mux.HandleFunc("POST /measurement", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		measurement := &measurement.Measurement{}
		if err := json.NewDecoder(r.Body).Decode(&measurement); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := brokerClient.Publish(config.BrokerSubject, measurement); err != nil {
			http.Error(w, "Failed to publish", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})

	stack := middleware.LoadMiddleware()
	server := http.Server{
		Addr:    config.GatewayAddr,
		Handler: stack(mux),
	}
	log.Println("Gateway starting on", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on %s: %v", server.Addr, err)
	}
}
