package measurement

import (
	"time"
)

type Measurement struct {
	ID        int       `json:"id"`
	DeviceId  string    `json:"deviceId"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type Device struct {
	ID        int       `json:"id"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"createdAt"`
}
