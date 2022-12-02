package main

import (
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rbobrovnikov/ejabberd-prometheus-metrics/core/config"
	"github.com/rbobrovnikov/ejabberd-prometheus-metrics/core/httprequest"
	"github.com/rbobrovnikov/ejabberd-prometheus-metrics/core/logger"
	"github.com/rbobrovnikov/ejabberd-prometheus-metrics/core/metrics"
)

func main() {
	logger.InitLogLevels(os.Stdout, os.Stdout, os.Stderr)
	logger.Info.Printf(config.ServiceName+" started at %s", config.ListenAddr)

	s, h, p, t := config.Config()
	logger.Info.Printf("Using ejabberd at %s://%s:%s with token: %s", s, h, p, t)
	http.Handle("/metrics", promhttp.Handler())
	ejabberdAPICheck()
	metrics.RecordMetrics(config.Config())
	logger.Info.Fatal(http.ListenAndServe(config.ListenAddr, nil))
}

// Check Ejabberd API availability
func ejabberdAPICheck() {
	s, h, p, t := config.Config()
	reqBodyJSONEmpty := `{}`
	ticker := time.NewTicker(5 * time.Second)
	for ejabberdAPIStatusCode := 0; ejabberdAPIStatusCode != 200; <-ticker.C {
		_, ejabberdAPIStatusCode = httprequest.EjabberAPICommonRequest(httprequest.HTTPBaseParams{Schema: s, Host: h, Port: p, Token: t, Endpoint: "status", ReqBody: reqBodyJSONEmpty})
		logger.Info.Printf("Status code: %d", ejabberdAPIStatusCode)
	}
	logger.Info.Printf("Ejabberd API is available. Ready to collect metrics!")
}
