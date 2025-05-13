package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"iotstarter/internal/model"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	monolithPort := flag.Bool("m", false, "Run in mode m (uses port 8080)")
	flag.Parse()

	port := "8081"
	if *monolithPort {
		port = "8080"
	}
	url := "http://localhost:" + port + "/api/measurements"
	log.Println("Sending data to: ", url)

	deviceIdStr := os.Getenv("TEST_DEVICE_ID")
	deviceId, err := strconv.Atoi(deviceIdStr)
	if err != nil {
		log.Fatalln("Incorrect device id format", deviceIdStr)
	}

	for {
		currentTime := time.Now().UTC()

		measurement := model.Measurement{
			DeviceId:  model.DeviceId(deviceId),
			Name:      "temperature",
			Value:     math.Sin(float64(currentTime.Unix()) / 10),
			Unit:      "C",
			Timestamp: currentTime,
		}
		err := sendMeasurement(url, measurement)
		if err != nil {
			log.Println("ERROR: ", err)
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func sendMeasurement(url string, measurement model.Measurement) error {
	jsonData, err := json.Marshal(measurement)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	apiKey := os.Getenv("TEST_DEVICE_API_KEY")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("Response Status:", resp.Status)
	return nil
}
