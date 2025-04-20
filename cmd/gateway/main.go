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
	brokerClient, err := broker.NewBrokerClient(cfg.BrokerUrl)
	if err != nil {
		log.Fatalln("ERROR: ", err.Error())
	}
	defer brokerClient.Close()

	handler := api.NewHandler().WithBroker(brokerClient)
	srv := api.NewServer(cfg.Addr, handler)
	srv.Run("Gateway")

	// mux := http.NewServeMux()

	// mux.HandleFunc("POST /measurement", func(w http.ResponseWriter, r *http.Request) {
	// 	defer r.Body.Close()
	// 	measurement := &measurement.Measurement{}
	// 	if err := json.NewDecoder(r.Body).Decode(&measurement); err != nil {
	// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	// 		return
	// 	}

	// 	if err := brokerClient.Publish(config.BrokerMeasurementSubject, measurement); err != nil {
	// 		http.Error(w, "Failed to publish", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.WriteHeader(http.StatusAccepted)
	// })

	// stack := middleware.LoadMiddleware()
	// server := http.Server{
	// 	Addr:    cfg.Addr,
	// 	Handler: stack(mux),
	// }
	// log.Println("Gateway starting on", server.Addr)
	// if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 	log.Fatalf("could not listen on %s: %v", server.Addr, err)
	// }
}
