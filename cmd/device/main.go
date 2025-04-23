package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"iotstarter/internal/logging"
	"iotstarter/internal/model"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	logging.SetUp()
	monolithPort := flag.Bool("m", false, "Run in mode m (uses port 8080)")
	flag.Parse()

	port := "8081"
	if *monolithPort {
		port = "8080"
	}
	url := "http://localhost:" + port + "/api/measurements"
	log.Println("Sending data to: ", url)

	for {
		measurement := model.Measurement{
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

func sendMeasurement(url string, measurement model.Measurement) error {
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
