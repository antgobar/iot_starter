package consumer

import (
	"context"
	"iotstarter/internal/model"
	"log"
	"time"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) Run() {
	err := h.svc.sub.Subscribe(context.Background(), h.svc.subject, h.saveMeasurement)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Transformer listening on subject: %s", h.svc.subject)
}

func (h *Handler) saveMeasurement(m *model.Measurement) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()
	err := h.svc.StoreMeasurement(ctx, m)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Stored measurement under id", m.ID, "for device id", m.DeviceId)
}
