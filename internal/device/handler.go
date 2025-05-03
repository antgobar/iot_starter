package device

import (
	"context"
	"encoding/json"
	"iotstarter/internal/auth"
	"iotstarter/internal/model"
	"iotstarter/internal/presentation"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	svc *Service
	p   presentation.Presenter
}

func NewHandler(svc *Service, p presentation.Presenter) *Handler {
	return &Handler{svc: svc, p: p}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /devices", h.register)
	mux.HandleFunc("GET /devices", h.list)
	mux.HandleFunc("GET /devices/{id}", h.getById)
	mux.HandleFunc("PATCH /devices/{id}/reauth", h.reauth)
	mux.HandleFunc("DELETE /devices/{id}", h.delete)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
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
	data := struct {
		Device *model.Device
	}{
		Device: device,
	}
	if err := h.p.Present(w, r, "device_row", data); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "resource error", http.StatusInternalServerError)
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(device)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	devices, err := h.svc.List(ctx, user.ID)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error retrieving devices", http.StatusInternalServerError)
		return
	}

	log.Println("devices found", len(devices))

	data := struct {
		User    *model.User
		Devices []*model.Device
	}{
		User:    user,
		Devices: devices,
	}
	if err := h.p.Present(w, r, "devices", data); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "resource error", http.StatusInternalServerError)
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(devices)
}

func (h *Handler) getById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
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
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

func (h *Handler) reauth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
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

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}
	deviceIdStr := r.PathValue("id")
	deviceId, err := strconv.Atoi(deviceIdStr)
	if err != nil {
		http.Error(w, "Invalid device Id", http.StatusBadRequest)
		return
	}

	dId := model.DeviceId(deviceId)

	if err := h.svc.DeleteDevice(ctx, user.ID, dId); err != nil {
		log.Println("ERROR:", err)
		http.Error(w, "Unable to offboard device with id: "+deviceIdStr, http.StatusUnauthorized)
		return
	}
}
