package transformer

import (
	"context"
	"iotstarter/internal/model"
	"iotstarter/internal/store"
	"log"
	"time"
)

type Handler struct {
	store *store.Store
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) SaveMeasurement(m *model.Measurement) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()
	err := h.store.SaveMeasurement(ctx, m)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Stored measurement under id", m.ID, "for device id", m.DeviceId)
}
