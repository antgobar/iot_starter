package measurement

import (
	"context"
	"encoding/json"
	"iotstarter/internal/auth"
	"iotstarter/internal/model"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/devices/{id}/measurements", h.getMeasurements)
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

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "error getting user from context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	deviceId := model.DeviceId(params.deviceId)
	measurements, err := h.svc.GetMeasurements(ctx, user.ID, deviceId, params.begin, params.end)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error retrieving measurements", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurements)
}
