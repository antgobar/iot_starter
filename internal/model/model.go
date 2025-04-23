package model

import (
	"time"
)

type Measurement struct {
	ID        int       `json:"id"`
	DeviceId  int       `json:"deviceId"`
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
}

type Device struct {
	ID        int       `json:"id"`
	UserId    int       `json:"userId"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"createdAt"`
	ApiKey    string    `json:"apiKey"`
}

type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	Active         bool      `json:"active"`
}
