package consumer

import (
	"context"
	"encoding/json"
	"log"

	"iotstarter/internal/measurement"
	"iotstarter/internal/model"
)

type Service struct {
	sub   Subscriber
	repo  measurement.Repository
	topic string
}

func NewService(sub Subscriber, repo measurement.Repository, topic string) *Service {
	return &Service{
		sub:   sub,
		repo:  repo,
		topic: topic,
	}
}

func (s *Service) handleMessage(ctx context.Context, payload []byte) {
	var m *model.Measurement
	if err := json.Unmarshal(payload, &m); err != nil {
		log.Printf("consumer: failed to unmarshal measurement: %v", err)
		return
	}

	if err := s.repo.Create(ctx, m); err != nil {
		log.Printf("consumer: failed to store measurement: %v", err)
		return
	}

	log.Printf("consumer: successfully stored measurement ID=%v", m.ID)
}
