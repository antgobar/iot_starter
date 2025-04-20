package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (h Handler) getDeviceMeasurements(w http.ResponseWriter, r *http.Request) {
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

func (h Handler) registerDevice(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	location := r.FormValue("location")
	if location == "" {
		http.Error(w, "Location not provided", http.StatusBadRequest)
		return
	}

	err := h.store.RegisterDevice(ctx, location)
	if err != nil {
		log.Println("ERROR:" + err.Error())
		http.Error(w, "Error registering device", http.StatusInternalServerError)
		return
	}
}

func (h Handler) getDevices(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	devices, err := h.store.GetDevices(ctx)
	if err != nil {
		log.Println("ERROR:" + err.Error())
		http.Error(w, "Error retrieving devices", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}
