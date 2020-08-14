package config

import (
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// ServiceName sets name of service
const ServiceName = "ejabberd-metrics-exporter"

// ListenAddr sets address to listen
var ListenAddr string = ":9334"

// EjabberdConfStuct configuration parameters
type EjabberdConfStuct struct {
	APIHost      string `required:"true" split_words:"true"`
	APIPort      string `default:"5443" split_words:"true"`
	APIUrlSchema string `default:"http" split_words:"true"`
	APIToken     string `required:"true" split_words:"true"`
}

// EjabberdCfg takes EjabberdConfStuct structure
var EjabberdCfg EjabberdConfStuct

// Config is for service configuration by environment variables
func Config() (string, string, string, string) {
	err := envconfig.Process(strings.Replace(ServiceName, "-", "_", -1), &EjabberdCfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	schema := EjabberdCfg.APIUrlSchema
	host := EjabberdCfg.APIHost
	port := EjabberdCfg.APIPort
	token := "Basic " + EjabberdCfg.APIToken

	return schema, host, port, token
}
