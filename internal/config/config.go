package config

import (
	"errors"
	"os"
)

func LoadGatewayConfig() (*GatewayConfig, error) {
	brokerUrl, err := loadEnv("BROKER_URL")
	if err != nil {
		return nil, err
	}
	brokerSubject, err := loadEnv("BROKER_SUBJECT")
	if err != nil {
		return nil, err
	}
	gatewayAddr, err := loadEnv("GATEWAY_ADDR")
	if err != nil {
		return nil, err
	}
	return &GatewayConfig{
		BrokerUrl:     brokerUrl,
		BrokerSubject: brokerSubject,
		GatewayAddr:   gatewayAddr}, nil
}

func LoadTransformerConfig() (*TransformerConfig, error) {
	brokerUrl, err := loadEnv("BROKER_URL")
	if err != nil {
		return nil, err
	}
	brokerSubject, err := loadEnv("BROKER_SUBJECT")
	if err != nil {
		return nil, err
	}
	dbUrl, err := loadEnv("DATABASE_URL")
	if err != nil {
		return nil, err
	}

	return &TransformerConfig{
		BrokerUrl:     brokerUrl,
		BrokerSubject: brokerSubject,
		DatabaseUrl:   dbUrl,
	}, nil
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
	GatewayAddr   string
}

type TransformerConfig struct {
	BrokerUrl     string
	BrokerSubject string
	DatabaseUrl   string
}
