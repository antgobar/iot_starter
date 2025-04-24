package consumer

import (
	"context"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/model"
	"log"
	"time"
)

type Consumer struct {
	subject string
	handler broker.MeasurementHandler
}

func newConsumer(subject string, handler broker.MeasurementHandler) Consumer {
	return Consumer{subject, handler}
}

func (h *Handler) registerConsumers() {
	h.consumers = nil
	var consumers []Consumer = []Consumer{
		newConsumer(config.BrokerMeasurementSubject, h.saveMeasurement),
	}
	h.consumers = append(h.consumers, consumers...)
}

func (h *Handler) saveMeasurement(m *model.Measurement) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()
	m, err := h.store.SaveMeasurement(ctx, m)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Stored measurement under id", m.ID, "for device id", m.DeviceId)
}
