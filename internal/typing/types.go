package typing

import "iotstarter/internal/model"

type MeasurementHandler func(msg *model.Measurement)
