package httprequest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rbobrovnikov/ejabberd-prometheus-metrics/core/logger"
)

// HTTPBaseParams defines basic http parameters
type HTTPBaseParams struct {
	Schema             string
	Host               string
	Port               string
	Token              string
	Endpoint           string
	ReqBody            string
	EjabberdMetricName string
}

// EjabberAPICommonRequest to Ejabberd HTTP API
func EjabberAPICommonRequest(p HTTPBaseParams) (float64, int) {

	url := p.Schema + "://" + p.Host + ":" + p.Port + "/api/" + p.Endpoint

	data := []byte(p.ReqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		logger.Error.Print("Error reading request. ", err)
		return 0, 0
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", p.Token)

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error.Print("Error reading response. ", err)
		return 0, 0
	}
	defer resp.Body.Close()

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	if p.EjabberdMetricName != "" {
		ejabberdMetricValue := (body[p.EjabberdMetricName].(float64))
		return ejabberdMetricValue, resp.StatusCode
	}
	return 0, resp.StatusCode
}
