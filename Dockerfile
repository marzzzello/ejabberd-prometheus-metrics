#
# Build arguments.
#
ARG GIT_TAG=unknown
ARG GIT_BRANCH=unknown
ARG GIT_COMMIT=unknown
ARG GO_SERVICE_IMPORT_PATH=unknown
ARG PATH_GO_SOURCES=/go/src/$GO_SERVICE_IMPORT_PATH


#
# Build Go binary inside base container.
#
FROM golang:1.14 as base_go_docker_image
# Env variables.
ARG GIT_TAG
ARG GIT_BRANCH
ARG GIT_COMMIT

ENV GIT_TAG=$GIT_TAG
ENV GIT_BRANCH=$GIT_BRANCH
ENV GIT_COMMIT=$GIT_COMMIT
# Create sources directory inside the container and copy project files.
RUN mkdir -p $GO_SERVICE_IMPORT_PATH/
WORKDIR $PATH_GO_SOURCES
RUN ls -la $PATH_GO_SOURCES
RUN ls -la $GO_SERVICE_IMPORT_PATH
COPY . $PATH_GO_SOURCES
# Build
RUN make build


#
# Destination container.
#
FROM scratch
LABEL key="Ruslan Bobrovnikov <ruslan.bobrovnikov@gmail.com>"
# Container arguments.
ARG GIT_TAG
ARG GIT_BRANCH
ARG GIT_COMMIT
# Container labels.
LABEL branch=$GIT_BRANCH
LABEL commit=$GIT_COMMIT
# Copy certificates and binary into the destination docker image.
COPY --from=base_go_docker_image /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base_go_docker_image /etc/passwd /etc/passwd
COPY --from=base_go_docker_image $PATH_GO_SOURCES/ejabberd-prometheus-metrics /ejabberd-prometheus-metrics
# Container settings.
ENV PORT 8080
EXPOSE 8080
USER nobody
ENTRYPOINT ["/ejabberd-prometheus-metrics"]
