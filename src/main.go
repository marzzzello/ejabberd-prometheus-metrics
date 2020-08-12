package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const serviceName = "ejabberd-metrics-exporter"

// EjabberdConf configuration parameters
type EjabberdConf struct {
	APIHost      string `required:"true" split_words:"true"`
	APIPort      string `default:"5443" split_words:"true"`
	APIUrlSchema string `default:"https" split_words:"true"`
	APIToken     string `required:"true" split_words:"true"`
}

// Define metrics
var (
	ejabberdConnectedUsersNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ejabberd",
			Name:      "connected_users_number",
			Help:      "The number of established sessions",
		})

	ejabberdIncommingS2SNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ejabberd",
			Name:      "incoming_s2s_number",
			Help:      "The number of incoming s2s connections on the node",
		})

	ejabberdOutgoingS2SNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ejabberd",
			Name:      "outgoing_s2s_number",
			Help:      "The number of outgoing s2s connections on the node",
		})

	ejabberdRegisteredUsers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ejabberd",
			Name:      "stats_registered_users",
			Help:      "The number of registered users",
		})

	ejabberdOnlineUsers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ejabberd",
			Name:      "stats_online_users",
			Help:      "The number of online users",
		})

	// ejabberdOnlineUsers = prometheus.NewGauge(
	// 	prometheus.GaugeOpts{
	// 		Namespace: "ejabberd",
	// 		Name:      "stats_online_users",
	// 		Help:      "The number of online users",
	// 	})
)

var listenAddr string = ":9334"

// Record metrics values
func recordMetrics(schema string, host string, port string, token string) {
	reqBodyJSONEmpty := `{}`
	go func() {
		for {
			ejabberdConnectedUsersNumber.Set(httpRequest(schema, host, port, token, "connected_users_number", reqBodyJSONEmpty, "num_sessions"))
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()

	go func() {
		for {
			ejabberdIncommingS2SNumber.Set(httpRequest(schema, host, port, token, "incoming_s2s_number", reqBodyJSONEmpty, "s2s_incoming"))
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()

	go func() {
		for {
			ejabberdOutgoingS2SNumber.Set(httpRequest(schema, host, port, token, "outgoing_s2s_number", reqBodyJSONEmpty, "s2s_outgoing"))
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()

	go func() {
		reqBodyJSON := `{"name": "registeredusers"}`
		for {
			ejabberdRegisteredUsers.Set(httpRequest(schema, host, port, token, "stats", reqBodyJSON, "stat"))
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()

	go func() {
		reqBodyJSON := `{"name": "onlineusers"}`
		for {
			ejabberdOnlineUsers.Set(httpRequest(schema, host, port, token, "stats", reqBodyJSON, "stat"))
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()
}

// HTTP requests to Ejabberd HTTP API
func httpRequest(schema string, host string, port string, token string, endpoint string, reqBody string, ejabberdMetricName string) float64 {

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

// Register metrics
func init() {
	prometheus.MustRegister(ejabberdConnectedUsersNumber)
	prometheus.MustRegister(ejabberdIncommingS2SNumber)
	prometheus.MustRegister(ejabberdOutgoingS2SNumber)
	prometheus.MustRegister(ejabberdRegisteredUsers)
	prometheus.MustRegister(ejabberdOnlineUsers)
}

func main() {
	var ejabberdCfg EjabberdConf
	err := envconfig.Process(strings.Replace(serviceName, "-", "_", -1), &ejabberdCfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	schema := ejabberdCfg.APIUrlSchema
	host := ejabberdCfg.APIHost
	port := ejabberdCfg.APIPort
	token := "Basic " + ejabberdCfg.APIToken

	recordMetrics(schema, host, port, token)

	http.Handle("/metrics", promhttp.Handler())
	log.Printf(serviceName+" started at %s\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))

}
