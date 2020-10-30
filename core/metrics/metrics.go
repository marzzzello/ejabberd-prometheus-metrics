package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rbobrovnikov/ejabberd-prometheus-metrics/core/httprequest"
)

type metricStruct struct {
	MetricName, MetricNamespace, MetricHelp, MetricEjabberdAPIEndpoint, MetricEjabberdAPIReqBody, MetricEjabberdAPIMetricSourceKey string
}

// Define metrics
var (
	ejabberdConnectedUsersNumber = metricStruct{
		MetricName:                       "connected_users_number",
		MetricHelp:                       "The number of established sessions",
		MetricEjabberdAPIEndpoint:        "connected_users_number",
		MetricEjabberdAPIMetricSourceKey: "num_sessions",
	}

	ejabberdIncommingS2SNumber = metricStruct{
		MetricName:                       "incoming_s2s_number",
		MetricHelp:                       "The number of incoming s2s connections on the node",
		MetricEjabberdAPIEndpoint:        "incoming_s2s_number",
		MetricEjabberdAPIMetricSourceKey: "s2s_incoming",
	}

	ejabberdOutgoingS2SNumber = metricStruct{
		MetricName:                       "outgoing_s2s_number",
		MetricHelp:                       "The number of outgoing s2s connections on the node",
		MetricEjabberdAPIEndpoint:        "outgoing_s2s_number",
		MetricEjabberdAPIMetricSourceKey: "s2s_outgoing",
	}

	ejabberdRegisteredUsers = metricStruct{
		MetricName:                       "stats_registered_users",
		MetricHelp:                       "The number of registered users",
		MetricEjabberdAPIEndpoint:        "stats",
		MetricEjabberdAPIMetricSourceKey: "stat",
		MetricEjabberdAPIReqBody:         `{"name": "registeredusers"}`,
	}

	ejabberdOnlineUsers = metricStruct{
		MetricName:                       "stats_online_users",
		MetricHelp:                       "The number of online users total",
		MetricEjabberdAPIEndpoint:        "stats",
		MetricEjabberdAPIMetricSourceKey: "stat",
		MetricEjabberdAPIReqBody:         `{"name": "onlineusers"}`,
	}

	ejabberdOnlineUsersNode = metricStruct{
		MetricName:                       "stats_online_users_node",
		MetricHelp:                       "The number of online users on the node",
		MetricEjabberdAPIEndpoint:        "stats",
		MetricEjabberdAPIMetricSourceKey: "stat",
		MetricEjabberdAPIReqBody:         `{"name": "onlineusersnode"}`,
	}

	ejabberdUptimeSeconds = metricStruct{
		MetricName:                       "stats_uptimeseconds",
		MetricHelp:                       "Uptime seconds",
		MetricEjabberdAPIEndpoint:        "stats",
		MetricEjabberdAPIMetricSourceKey: "stat",
		MetricEjabberdAPIReqBody:         `{"name": "uptimeseconds"}`,
	}

	ejabberdProcesses = metricStruct{
		MetricName:                       "stats_processes",
		MetricHelp:                       "The number of processes",
		MetricEjabberdAPIEndpoint:        "stats",
		MetricEjabberdAPIMetricSourceKey: "stat",
		MetricEjabberdAPIReqBody:         `{"name": "processes"}`,
	}

	ejabberdClusterNodesCount = metricStruct{
		MetricName:                "cluster_nodes_count",
		MetricHelp:                "The number of cluster nodes",
		MetricEjabberdAPIEndpoint: "list_cluster",
	}

	ejabberdMUCNumber = metricStruct{
		MetricName:                "muc_online_rooms",
		MetricHelp:                "Number of existing rooms (MUC)",
		MetricEjabberdAPIEndpoint: "muc_online_rooms",
		MetricEjabberdAPIReqBody:  `{"service": "global"}`,
	}

	ejabberdMUCRoomsEmptyNumber = metricStruct{
		MetricName:                "muc_rooms_empty_number",
		MetricHelp:                "Number of the rooms (MUC) that have no messages in archive",
		MetricEjabberdAPIEndpoint: "rooms_empty_list",
		MetricEjabberdAPIReqBody:  `{"service": "global"}`,
	}

	ejabberdMUCRoomsUnusedNumber30days = metricStruct{
		MetricName:                "muc_rooms_unused_30days_number",
		MetricHelp:                "Number of the rooms (MUC) that are unused for 30 days in the service",
		MetricEjabberdAPIEndpoint: "rooms_unused_list",
		MetricEjabberdAPIReqBody:  `{"service": "global", "days": 30}`,
	}

	ejabberdMUCRoomsUnusedNumber90days = metricStruct{
		MetricName:                "muc_rooms_unused_90days_number",
		MetricHelp:                "Number of the rooms (MUC) that are unused for 90 days in the service",
		MetricEjabberdAPIEndpoint: "rooms_unused_list",
		MetricEjabberdAPIReqBody:  `{"service": "global", "days": 90}`,
	}

	ejabberdMUCRoomsUnusedNumber180days = metricStruct{
		MetricName:                "muc_rooms_unused_180days_number",
		MetricHelp:                "Number of the rooms (MUC) that are unused for 180 days in the service",
		MetricEjabberdAPIEndpoint: "rooms_unused_list",
		MetricEjabberdAPIReqBody:  `{"service": "global", "days": 180}`,
	}

	ejabberdMUCRoomsUnusedNumber360days = metricStruct{
		MetricName:                "muc_rooms_unused_360days_number",
		MetricHelp:                "Number of the rooms (MUC) that are unused for 360 days in the service",
		MetricEjabberdAPIEndpoint: "rooms_unused_list",
		MetricEjabberdAPIReqBody:  `{"service": "global", "days": 360}`,
	}
)

// RecordMetrics generates metrics
func RecordMetrics(schema string, host string, port string, token string) {

	metricsList := []metricStruct{
		ejabberdConnectedUsersNumber,
		ejabberdIncommingS2SNumber,
		ejabberdOutgoingS2SNumber,
		ejabberdRegisteredUsers,
		ejabberdOnlineUsers,
		ejabberdOnlineUsersNode,
		ejabberdUptimeSeconds,
		ejabberdProcesses,
		ejabberdClusterNodesCount,
		ejabberdMUCNumber,
		ejabberdMUCRoomsEmptyNumber,
		ejabberdMUCRoomsUnusedNumber30days,
		ejabberdMUCRoomsUnusedNumber90days,
		ejabberdMUCRoomsUnusedNumber180days,
		ejabberdMUCRoomsUnusedNumber360days,
	}

	for _, m := range metricsList {
		// Set metrics default values
		if m.MetricNamespace == "" {
			m.MetricNamespace = "ejabberd"
		}
		if m.MetricEjabberdAPIReqBody == "" {
			m.MetricEjabberdAPIReqBody = `{}`
		}

		newMetricRoutine(schema, host, port, token, m)
	}
}

func newMetricRoutine(schema string, host string, port string, token string, m metricStruct) {
	scrapeInterval := (time.Duration(5) * time.Second)
	ticker := time.NewTicker(scrapeInterval)

	prometheusMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: m.MetricNamespace,
			Name:      m.MetricName,
			Help:      m.MetricHelp,
		})

	go func() {
		for range ticker.C {
			ejabberdMetricValue, _ := httprequest.EjabberAPICommonRequest(httprequest.HTTPBaseParams{Schema: schema, Host: host, Port: port, Token: token, Endpoint: m.MetricEjabberdAPIEndpoint, ReqBody: m.MetricEjabberdAPIReqBody, EjabberdAPIMetricSourceKey: m.MetricEjabberdAPIMetricSourceKey})
			prometheusMetric.Set(ejabberdMetricValue)
		}
	}()
	prometheus.MustRegister(prometheusMetric)
}
