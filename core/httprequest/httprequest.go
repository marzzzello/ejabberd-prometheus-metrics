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
	Schema, Host, Port, Token, Endpoint, ReqBody, EjabberdAPIMetricSourceKey string
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

	if p.EjabberdAPIMetricSourceKey == "" {
		var respArr []string
		json.NewDecoder(resp.Body).Decode(&respArr)
		ejabberdMetricValue := len(respArr)
		return float64(ejabberdMetricValue), resp.StatusCode
	}

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	ejabberdMetricValue := (body[p.EjabberdAPIMetricSourceKey].(float64))
	return ejabberdMetricValue, resp.StatusCode

}
