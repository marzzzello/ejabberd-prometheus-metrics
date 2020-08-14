package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rbobrovnikov/ejabberd-prometheus-metrics/core/config"
	"github.com/rbobrovnikov/ejabberd-prometheus-metrics/core/logger"
	"github.com/rbobrovnikov/ejabberd-prometheus-metrics/core/metrics"
)

func main() {
	logger.LoggerInit(os.Stdout, os.Stdout, os.Stderr)

	metrics.RegisterMetrics()
	metrics.RecordMetrics(config.Config())

	http.Handle("/metrics", promhttp.Handler())

	logger.Info.Printf(config.ServiceName+" started at %s\n\nBuild info:\nDate: %s\nTag: %s\nBranch: %s\nCommit: %s\n", config.ListenAddr, logger.BuildDate, logger.BuildTag, logger.BuildBranch, logger.BuildCommit)
	log.Fatal(http.ListenAndServe(config.ListenAddr, nil))

}
