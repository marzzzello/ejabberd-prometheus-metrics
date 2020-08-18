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
	logger.Info.Printf(config.ServiceName+" started at %s\nBuild info: [ Date: %s | Tag: %s | Branch: %s | Commit: %s ]", config.ListenAddr, logger.BuildDate, logger.BuildTag, logger.BuildBranch, logger.BuildCommit)

	metrics.RegisterMetrics()
	http.Handle("/metrics", promhttp.Handler())
	ejabberdAPICheck()
	logger.Info.Fatal(http.ListenAndServe(config.ListenAddr, nil))
}

// Check Ejabberd API availability
func ejabberdAPICheck() {
	for ejabberdAPIStatusCode := 0; ejabberdAPIStatusCode != 200; {
		s, h, p, t := config.Config()
		reqBodyJSONEmpty := `{}`
		_, ejabberdAPIStatusCode = httprequest.EjabberAPICommonRequest(httprequest.HTTPBaseParams{Schema: s, Host: h, Port: p, Token: t, Endpoint: "status", ReqBody: reqBodyJSONEmpty})
		time.Sleep(5 * time.Second)
	}
	metrics.RecordMetrics(config.Config())
	logger.Info.Printf("Ejabberd API is available. Ready to collect metrics!")
}
