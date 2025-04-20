package main

import (
	"bytes"
	"encoding/json"
	"iotstarter/internal/config"
	"iotstarter/internal/measurement"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	gatewayUrl := config.MustLoadEnv("GATEWAY_URL")

	url := gatewayUrl + "/measurements"
	log.Println("Sending data to: ", url)
	for {
		measurement := measurement.Measurement{
			DeviceId:  1,
			Name:      "temperature",
			Value:     rand.Float64() * 10,
			Unit:      "C",
			Timestamp: time.Now().UTC(),
		}
		err := sendMeasurement(url, measurement)
		if err != nil {
			log.Println("ERROR: ", err)
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func sendMeasurement(url string, measurement measurement.Measurement) error {
	jsonData, err := json.Marshal(measurement)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("Response Status:", resp.Status)
	return nil
}
