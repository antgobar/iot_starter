package broker

import "iotstarter/internal/model"

type MeasurementHandler func(msg *model.Measurement)

type Broker interface {
	Publish(subject string, msg *model.Measurement) error
	Subscribe(subject string, handler MeasurementHandler) error
}
