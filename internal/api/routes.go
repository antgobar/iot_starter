package api

import (
	"context"
	"encoding/json"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/model"
	"iotstarter/internal/store"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	store  *store.Store
	broker broker.Broker
}

func NewHandler(store *store.Store, broker broker.Broker) *Handler {
	return &Handler{store, broker}
}

func (h *Handler) registerUserRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	if h.store != nil {
		mux.HandleFunc("POST /devices", h.registerDevice)
		mux.HandleFunc("GET /devices", h.getDevices)
		mux.HandleFunc("GET /devices/{id}/measurements", h.getDeviceMeasurements)
	}

	if h.broker != nil {
		mux.HandleFunc("POST /measurements", h.saveMeasurement)
	}
	return mux
}

func (h *Handler) getDeviceMeasurements(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	params, err := getMeasurementsQueryParams(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	measurements, err := h.store.GetDeviceMeasurements(ctx, params.deviceId, params.start, params.end)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error retrieving measurements", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurements)
}

func (h *Handler) registerDevice(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	location := r.FormValue("location")
	if location == "" {
		http.Error(w, "location required", http.StatusBadRequest)
		return
	}

	err := h.store.RegisterDevice(ctx, location)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error registering device", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getDevices(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	devices, err := h.store.GetDevices(ctx)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error retrieving devices", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func (h *Handler) saveMeasurement(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	measurement := &model.Measurement{}
	if err := json.NewDecoder(r.Body).Decode(&measurement); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.broker.Publish(config.BrokerMeasurementSubject, measurement); err != nil {
		http.Error(w, "Failed to publish", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
