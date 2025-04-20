package broker

import (
	"encoding/json"
	"iotstarter/internal/measurement"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type BrokerClient struct {
	Connection *nats.Conn
}

func NewBrokerClient(connectionString string) (*BrokerClient, error) {
	nc, err := nats.Connect(connectionString)
	if err != nil {
		return nil, err
	}
	return &BrokerClient{Connection: nc}, nil
}

func (b BrokerClient) Publish(subject string, measurement *measurement.Measurement) error {
	if measurement.Timestamp.IsZero() {
		measurement.Timestamp = time.Now().UTC()
		log.Println("Timestamp not provided: calculated on publish")
	}
	data, err := json.Marshal(measurement)
	if err != nil {
		return err
	}

	if err := b.Connection.Publish(subject, data); err != nil {
		return err
	}

	return nil
}

func (b BrokerClient) Close() {
	b.Connection.Close()
}

func (b BrokerClient) Subscribe(subject string, handler MeasurementHandler) error {
	processMessage := func(msg *nats.Msg) {
		var measurement measurement.Measurement
		if err := json.Unmarshal(msg.Data, &measurement); err != nil {
			log.Printf("Error decoding message: %v", err)
			return
		}
		handler(&measurement)
	}
	_, err := b.Connection.Subscribe(subject, processMessage)
	if err != nil {
		return err
	}
	return nil
}
