package consumer

import (
	"context"
	"iotstarter/internal/config"
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
	err := h.svc.sub.Subscribe(context.TODO(), config.BrokerMeasurementSubject, h.saveMeasurement)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Transformer listening on subject: %s", config.BrokerMeasurementSubject)
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
