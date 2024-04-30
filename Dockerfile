FROM golang:1.22-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go .
COPY core ./core

RUN go build -o /ejabberd-prometheus-metrics

##
## Deploy
##
FROM ubuntu
RUN apt update && apt install -y net-tools

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /ejabberd-prometheus-metrics /ejabberd-prometheus-metrics

EXPOSE 9100

USER nobody

ENTRYPOINT ["/ejabberd-prometheus-metrics"]
