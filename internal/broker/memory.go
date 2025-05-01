package broker

import (
	"context"
	"iotstarter/internal/model"
	"iotstarter/internal/typing"
	"sync"
)

type memoryBroker struct {
	publishCh   chan *publishRequest
	subscribeCh chan *subscribeRequest
	handlers    map[string][]typing.MeasurementHandler
	mu          sync.Mutex
}

type publishRequest struct {
	subject string
	msg     *model.Measurement
	ack     chan error
}

type subscribeRequest struct {
	subject string
	handler typing.MeasurementHandler
	ack     chan error
}

func NewMemoryBroker() Broker {
	b := &memoryBroker{
		publishCh:   make(chan *publishRequest),
		subscribeCh: make(chan *subscribeRequest),
		handlers:    make(map[string][]typing.MeasurementHandler),
	}
	go b.run()
	return b
}

func (b *memoryBroker) run() {
	for {
		select {
		case sub := <-b.subscribeCh:
			b.mu.Lock()
			b.handlers[sub.subject] = append(b.handlers[sub.subject], sub.handler)
			b.mu.Unlock()
			sub.ack <- nil

		case pub := <-b.publishCh:
			b.mu.Lock()
			handlers := b.handlers[pub.subject]
			b.mu.Unlock()

			for _, handler := range handlers {
				go handler(pub.msg)
			}
			pub.ack <- nil
		}
	}
}

func (b *memoryBroker) Subscribe(ctx context.Context, subject string, handler typing.MeasurementHandler) error {
	ack := make(chan error)
	b.subscribeCh <- &subscribeRequest{subject: subject, handler: handler, ack: ack}
	return <-ack
}

func (b *memoryBroker) Publish(ctx context.Context, subject string, msg *model.Measurement) error {
	ack := make(chan error)
	b.publishCh <- &publishRequest{subject: subject, msg: msg, ack: ack}
	return <-ack
}
