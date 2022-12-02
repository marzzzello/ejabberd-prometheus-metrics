# ejabberd-prometheus-metrics
Service that collects Ejabberd stats and exports them in Prometheus format

Service takes configuration from environment variables:
```
EJABBERD_METRICS_EXPORTER_API_URL_SCHEMA (optional, default: http)
EJABBERD_METRICS_EXPORTER_API_HOST
EJABBERD_METRICS_EXPORTER_API_PORT       (optional, default: 5443)
EJABBERD_METRICS_EXPORTER_API_TOKEN
```

Metrics exposed via endpoint `/metrics` and port 9100.

docker-compose.yml:
```yml
services:
  ejabberd-exporter:
    image: ghcr.io/marzzzello/ejabberd-prometheus-metrics:main
    container_name: ejabberd-exporter
    restart: unless-stopped
    networks:
      - your_ejabberd_network
      - your_metrics_network
    user: nobody # default user
    security_opt:
      - no-new-privileges=true
    cap_drop: # drop unneeded capabilities
      - ALL
    environment:
      # EJABBERD_METRICS_EXPORTER_API_URL_SCHEMA: https
      EJABBERD_METRICS_EXPORTER_API_HOST: ejabberd
      EJABBERD_METRICS_EXPORTER_API_PORT: 5080
      EJABBERD_METRICS_EXPORTER_API_TOKEN: ${EJABBERD_API_TOKEN}
```
