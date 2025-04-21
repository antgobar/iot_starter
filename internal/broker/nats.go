package broker

import (
	"encoding/json"
	"iotstarter/internal/model"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsBrokerClient struct {
	nc *nats.Conn
}

func NewNatsBrokerClient(connectionString string) (*NatsBrokerClient, error) {
	nc, err := nats.Connect(connectionString)
	if err != nil {
		return nil, err
	}
	return &NatsBrokerClient{nc: nc}, nil
}

func (b *NatsBrokerClient) Publish(subject string, measurement *model.Measurement) error {
	if measurement.Timestamp.IsZero() {
		measurement.Timestamp = time.Now().UTC()
		log.Println("Timestamp not provided: calculated on publish")
	}
	data, err := json.Marshal(measurement)
	if err != nil {
		return err
	}

	if err := b.nc.Publish(subject, data); err != nil {
		return err
	}

	return nil
}

func (b *NatsBrokerClient) Close() {
	err := b.nc.Drain()
	if err != nil {
		log.Println("Error draining NATS", err.Error())
	}
	b.nc.Close()
}

func (b *NatsBrokerClient) Subscribe(subject string, handler MeasurementHandler) error {
	processMessage := func(msg *nats.Msg) {
		var measurement model.Measurement
		if err := json.Unmarshal(msg.Data, &measurement); err != nil {
			log.Printf("Error decoding message: %v", err)
			return
		}
		go handler(&measurement)
	}
	_, err := b.nc.Subscribe(subject, processMessage)
	if err != nil {
		return err
	}
	return nil
}
