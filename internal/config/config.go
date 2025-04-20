package config

import (
	"errors"
	"os"
)

const BrokerMeasurementSubject = "measurement"

func LoadGatewayConfig() (*GatewayConfig, error) {
	brokerUrl, err := loadEnv("BROKER_URL")
	if err != nil {
		return nil, err
	}
	gatewayAddr, err := loadEnv("GATEWAY_ADDR")
	if err != nil {
		return nil, err
	}
	return &GatewayConfig{
		BrokerUrl:     brokerUrl,
		BrokerSubject: BrokerMeasurementSubject,
		Addr:          gatewayAddr}, nil
}

func LoadTransformerConfig() (*TransformerConfig, error) {
	brokerUrl, err := loadEnv("BROKER_URL")
	if err != nil {
		return nil, err
	}
	dbUrl, err := loadEnv("DATABASE_URL")
	if err != nil {
		return nil, err
	}

	return &TransformerConfig{
		BrokerUrl:     brokerUrl,
		BrokerSubject: BrokerMeasurementSubject,
		DatabaseUrl:   dbUrl,
	}, nil
}

func LoadDashboardConfig() (*DashboardConfig, error) {
	apiAddr, err := loadEnv("API_ADDR")
	if err != nil {
		return nil, err
	}
	dbUrl, err := loadEnv("DATABASE_URL")
	if err != nil {
		return nil, err
	}
	return &DashboardConfig{Addr: apiAddr, DatabaseUrl: dbUrl}, nil
}

func loadEnv(envName string) (string, error) {
	env := os.Getenv(envName)
	if env == "" {
		return "", errors.New("missing environment variable " + envName)
	}
	return env, nil
}

type GatewayConfig struct {
	BrokerUrl     string
	BrokerSubject string
	Addr          string
}

type TransformerConfig struct {
	BrokerUrl     string
	BrokerSubject string
	DatabaseUrl   string
}

type DashboardConfig struct {
	Addr        string
	DatabaseUrl string
}
