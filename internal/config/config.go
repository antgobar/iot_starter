package config

import (
	"log"
	"os"
)

const BrokerMeasurementSubject = "measurement"

func MustLoadEnv(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		log.Fatalf("missing environment variable: %s", envName)
	}
	return env
}

type RequestContextKey string

const UserKey RequestContextKey = "user"
