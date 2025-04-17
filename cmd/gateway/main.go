package main

import (
	"encoding/json"
	"errors"
	"iotstarter/internal/broker"
	"iotstarter/internal/measurement"
	"iotstarter/internal/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	config, err := GetConfig()
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	brokerClient, err := broker.NewBrokerClient(config.BrokerAddr)
	if err != nil {
		log.Println("ERROR: error connecting to broker client")
		return
	}
	defer brokerClient.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /measurement/schema", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		example := measurement.Measurement{
			DeviceId: "your-device-id",
			Name:     "your-measurement-name",
			Value:    "some-value",
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

		if err := brokerClient.Publish(config.Subject, measurement); err != nil {
			http.Error(w, "Failed to publish", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})

	stack := middleware.LoadMiddleware()
	server := http.Server{
		Addr:    config.ServerAddr,
		Handler: stack(mux),
	}
	log.Println("Server starting on", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on %s: %v", server.Addr, err)
	}
}

type Config struct {
	ServerAddr string
	Subject    string
	BrokerAddr string
}

func GetConfig() (*Config, error) {
	serverAddr, err := loadEnv("APP_ADDR")
	if err != nil {
		return nil, err
	}
	brokerSubject, err := loadEnv("BROKER_SUBJECT")
	if err != nil {
		return nil, err
	}
	brokerUrl, err := loadEnv("BROKER_URL")
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerAddr: serverAddr,
		Subject:    brokerSubject,
		BrokerAddr: brokerUrl,
	}, nil
}

func loadEnv(envName string) (string, error) {
	env := os.Getenv(envName)
	if env == "" {
		return "", errors.New("missing environment variable " + envName)
	}
	return env, nil
}
