package consumer

import (
	"context"
	"iotstarter/internal/config"
	"iotstarter/internal/model"
	"iotstarter/internal/typing"
	"log"
	"time"
)

type Handler struct {
	svc       *Service
	consumers []*Consumer
}

type Consumer struct {
	subject string
	handler typing.MeasurementHandler
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc, nil}
}



func (h *Handler) Run() {
	h.registerConsumers()
	for _, consumer := range h.consumers {
		err := h.svc.subscriber.Subscribe(context.TODO(), consumer.subject, consumer.handler)
		if err != nil {
			log.Fatalln(err.Error())
		}
		h.consumers = append(h.consumers, consumer)
	}
	log.Printf("Transformer listening on subject(s): %s", h.consumersSubjects())
	select {}
}

func (h *Handler) consumersSubjects() []string {
	subjects := make([]string, 0)
	for _, consumer := range h.consumers {
		subjects = append(subjects, consumer.subject)
	}
	return subjects
}

func (h *Handler) registerConsumers() {
	h.consumers = nil
	var consumers []*Consumer = []*Consumer{
		newConsumer(config.BrokerMeasurementSubject, h.saveMeasurement),
	}
	h.consumers = append(h.consumers, consumers...)
}

func newConsumer(subject string, handler typing.MeasurementHandler) *Consumer {
	return &Consumer{subject, handler}
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
