# ejabberd-prometheus-metrics
Service that collects Ejabberd stats and exports them in Prometheus format

Service takes configuration from environment variables:
```
EJABBERD_METRICS_EXPORTER_API_URL_SCHEMA (optional, default: http)
EJABBERD_METRICS_EXPORTER_API_HOST        
EJABBERD_METRICS_EXPORTER_API_PORT       (optional, default: 5443)
EJABBERD_METRICS_EXPORTER_API_TOKEN
```

Metrics exposed via endpoint `/metrics`
