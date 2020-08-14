package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rbobrovnikov/ejabberd-prometheus-metrics/core/http"
)

// Define metrics
var (
	EjabberdConnectedUsersNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ejabberd",
			Name:      "connected_users_number",
			Help:      "The number of established sessions",
		})

	EjabberdIncommingS2SNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ejabberd",
			Name:      "incoming_s2s_number",
			Help:      "The number of incoming s2s connections on the node",
		})

	EjabberdOutgoingS2SNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ejabberd",
			Name:      "outgoing_s2s_number",
			Help:      "The number of outgoing s2s connections on the node",
		})

	EjabberdRegisteredUsers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ejabberd",
			Name:      "stats_registered_users",
			Help:      "The number of registered users",
		})

	EjabberdOnlineUsers = prometheus.NewGauge(
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

// RecordMetrics generates metrics
func RecordMetrics(schema string, host string, port string, token string) {
	reqBodyJSONEmpty := `{}`
	scrapeInterval := (time.Duration(5) * time.Second)
	go func() {
		for {
			EjabberdConnectedUsersNumber.Set(http.HttpRequest(schema, host, port, token, "connected_users_number", reqBodyJSONEmpty, "num_sessions"))
			time.Sleep(scrapeInterval)
		}
	}()

	go func() {
		for {
			EjabberdIncommingS2SNumber.Set(http.HttpRequest(schema, host, port, token, "incoming_s2s_number", reqBodyJSONEmpty, "s2s_incoming"))
			time.Sleep(scrapeInterval)
		}
	}()

	go func() {
		for {
			EjabberdOutgoingS2SNumber.Set(http.HttpRequest(schema, host, port, token, "outgoing_s2s_number", reqBodyJSONEmpty, "s2s_outgoing"))
			time.Sleep(scrapeInterval)
		}
	}()

	go func() {
		reqBodyJSON := `{"name": "registeredusers"}`
		for {
			EjabberdRegisteredUsers.Set(http.HttpRequest(schema, host, port, token, "stats", reqBodyJSON, "stat"))
			time.Sleep(scrapeInterval)
		}
	}()

	go func() {
		reqBodyJSON := `{"name": "onlineusers"}`
		for {
			EjabberdOnlineUsers.Set(http.HttpRequest(schema, host, port, token, "stats", reqBodyJSON, "stat"))
			time.Sleep(scrapeInterval)
		}
	}()
}

// RegisterMetrics sets up configured metrics
func RegisterMetrics() {
	prometheus.MustRegister(EjabberdConnectedUsersNumber)
	prometheus.MustRegister(EjabberdIncommingS2SNumber)
	prometheus.MustRegister(EjabberdOutgoingS2SNumber)
	prometheus.MustRegister(EjabberdRegisteredUsers)
	prometheus.MustRegister(EjabberdOnlineUsers)
}
