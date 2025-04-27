package model

import (
	"time"
)

type UserId int

type User struct {
	ID             UserId    `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	Active         bool      `json:"active"`
}

type DeviceId int
type ApiKey string

type Device struct {
	ID        DeviceId  `json:"id"`
	UserId    UserId    `json:"userId"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"createdAt"`
	ApiKey    ApiKey    `json:"apiKey"`
}

type MeasurementID int

type Measurement struct {
	ID        MeasurementID `json:"id"`
	DeviceId  DeviceId      `json:"deviceId"`
	Name      string        `json:"name"`
	Value     float64       `json:"value"`
	Unit      string        `json:"unit"`
	Timestamp time.Time     `json:"timestamp"`
}

type SessionID int

type Session struct {
	ID        SessionID `json:"id"`
	UserId    UserId    `json:"userId"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}
