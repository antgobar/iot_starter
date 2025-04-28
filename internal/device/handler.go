package device

import (
	"context"
	"encoding/json"
	"iotstarter/internal/auth"
	"iotstarter/internal/model"
	"iotstarter/internal/store"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/devices", h.register)
	mux.HandleFunc("GET /api/devices", h.list)
	mux.HandleFunc("GET /api/devices/{id}", h.getById)
	mux.HandleFunc("GET /api/devices/{id}/measurements", h.getMeasurements)
	mux.HandleFunc("PATCH /api/devices/{id}/reauth", h.reauth)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	user, err := auth.UserFromContext(r.Context())
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	location := r.FormValue("location")
	if location == "" {
		http.Error(w, "location required", http.StatusBadRequest)
		return
	}

	device, err := h.svc.Register(ctx, user.ID, location)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error registering device", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	user, err := auth.UserFromContext(r.Context())
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	devices, err := h.svc.List(ctx, user.ID)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error retrieving devices", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func (h *Handler) getById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	user, err := auth.UserFromContext(r.Context())
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	deviceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid device Id", http.StatusBadRequest)
		return
	}

	deviceIdModel := model.DeviceId(deviceId)
	device, err := h.svc.GetUserDeviceById(ctx, user.ID, deviceIdModel)
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

func (h *Handler) getMeasurements(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	params, err := getMeasurementsQueryParams(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user, err := auth.UserFromContext(r.Context())
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	deviceId := model.DeviceId(params.deviceId)

	measurements, err := h.svc.GetMeasurements(ctx, user.ID, deviceId, params.start, params.end)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error retrieving measurements", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurements)
}

func (h *Handler) reauth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	user, err := auth.UserFromContext(r.Context())
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	deviceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid device Id", http.StatusBadRequest)
		return
	}

	dId := model.DeviceId(deviceId)

	device, err := h.svc.Reauth(ctx, user.ID, dId)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error reauthing device", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}
