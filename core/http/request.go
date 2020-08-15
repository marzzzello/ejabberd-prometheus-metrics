package http

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// EjabberAPICommonRequest to Ejabberd HTTP API
func EjabberAPICommonRequest(schema, host, port, token, endpoint, reqBody, ejabberdMetricName string) float64 {

	url := schema + "://" + host + ":" + port + "/api/" + endpoint

	data := []byte(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	ejabberdMetricValue := (body[ejabberdMetricName].(float64))

	return ejabberdMetricValue
}
