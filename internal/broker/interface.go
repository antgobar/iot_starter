package broker

import "iotstarter/internal/measurement"

type MeasurementHandler func(msg *measurement.Measurement)

type Broker interface {
	Publish(subject string, msg *measurement.Measurement) error
	Subscribe(subject string, handler MeasurementHandler) error
}
