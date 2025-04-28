package gateway

import (
	"context"
	"encoding/json"
	"iotstarter/internal/model"
	"net/http"
	"time"
)

type Handler struct {
	svc     *Service
	subject string
}

func NewHandler(svc *Service, subject string) *Handler {
	return &Handler{svc: svc, subject: subject}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/devices", h.saveMeasurement)
}

func (h *Handler) saveMeasurement(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	measurement := &model.Measurement{}
	if err := json.NewDecoder(r.Body).Decode(&measurement); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.svc.Publish(ctx, h.subject, measurement); err != nil {
		http.Error(w, "Failed to publish", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
