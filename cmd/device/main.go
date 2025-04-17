package main

import (
	"bytes"
	"encoding/json"
	"iotstarter/internal/measurement"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func main() {
	serverAddr := os.Getenv("COLLECTOR_URL")
	if serverAddr == "" {
		log.Println("No server address provided")
		return
	}

	url := serverAddr + "/measurement"
	log.Println("Sending data to: ", url)
	deviceId := uuid.New().String()
	for {
		randomInt := rand.Intn(100)
		measurement := measurement.Measurement{
			DeviceId:  deviceId,
			Name:      "temperature",
			Value:     strconv.Itoa(randomInt),
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
