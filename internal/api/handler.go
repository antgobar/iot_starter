package api

import (
	"context"
	"encoding/json"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/model"
	"iotstarter/internal/store"
	"iotstarter/internal/views"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	store  *store.Store
	broker broker.Broker
	views  *views.Views
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) WithStore(store *store.Store) *Handler {
	h.store = store
	return h
}

func (h *Handler) WithBroker(broker broker.Broker) *Handler {
	h.broker = broker
	return h
}

func (h *Handler) WithViews(views *views.Views) *Handler {
	h.views = views
	return h
}

func (h *Handler) registerUserRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	if h.store != nil {
		mux.HandleFunc("POST /devices", h.registerDevice)
		mux.HandleFunc("GET /devices", h.getDevices)
		mux.HandleFunc("GET /devices/{id}", h.getDeviceById)
		mux.HandleFunc("GET /devices/{id}/measurements", h.getDeviceMeasurements)
		mux.HandleFunc("GET /measurements", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/devices", http.StatusFound)
		})
	}

	if h.broker != nil {
		mux.HandleFunc("POST /measurements", h.saveMeasurement)
	}

	if h.views != nil {
		mux.HandleFunc("GET /", h.getIndexPage)
	}

	return mux
}

func (h *Handler) getIndexPage(w http.ResponseWriter, r *http.Request) {
	view := h.views.IndexPage(w)
	view.Render(nil)
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

func (h *Handler) getDeviceById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	deviceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid device Id", http.StatusBadRequest)
		return
	}

	device, err := h.store.GetDeviceById(ctx, deviceId)
	if err == store.ErrDeviceNotFound {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error retrieving devices", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
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
