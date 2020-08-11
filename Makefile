# Virgil Makefile Version 1.0.0
.DEFAULT_GOAL := build

# Target list
.PHONY: init build clean
.PHONY: go_get go_build
.PHONY: docker_inspect docker_build_image
.PHONY: docker_registry_push



#
# Project-specific variables
#
# Service name. Used for binary name, docker-compose service name, etc...
SERVICE=ejabberd-prometheus-metrics
# Path to service entry point.
GO_PATH_SERVICE_MAIN=./src

#
# General variables
#
# Current build commit.
GO_SERVICE_BUILD_COMMIT=$(GIT_COMMIT)
# Current build branch.
GO_SERVICE_BUILD_BRANCH=$(GIT_BRANCH)
# Current build tag.
GO_SERVICE_BUILD_TAG=$(GIT_TAG)
# Build platform
BUILD_PLATFORM=$(shell go env GOOS)_$(shell go env GOARCH)
# Docker repository
DOCKER_REPO=rbobrovnikov/$(SERVICE)
# Path to Docker file
PATH_DOCKER_FILE=$(realpath Dockerfile)


#
# Build targets
#

# init target substitutes all the template variables in build files (like Dockerfile and docer-compose configs) with
# correct their correct values.
init:
	@echo '>>> Initting the project.'
	@sed -i'' "s/{% SERVICE_NAME %}/$(SERVICE)/g" $(PATH_DOCKER_FILE);

build: clean go_get go_build

clean:
	@echo ">>> Cleaning."
	@rm -rvf $(SERVICE)

#
# Go targets
#
go_get:
	@echo '>>> Getting go modules.'
	@go get -v -t -d ./...

go_build:
	@echo '>>> Building go binary.'
	@go build -o $(SERVICE) $(GO_PATH_SERVICE_MAIN);

#
# Docker targets
#
docker_inspect:
	@echo ">>> Inspecting docker container."
	docker inspect \
		-f '{{index .ContainerConfig.Labels "version"}}' \
		-f '{{index .ContainerConfig.Labels "branch"}}' \
		$(SERVICE)

docker_build_image:
	@echo ">>> Building docker image with service binary."
	docker build \
		-t $(SERVICE) \
		--build-arg GIT_TAG=$(GO_SERVICE_BUILD_TAG) \
		--build-arg GIT_BRANCH=$(GO_SERVICE_BUILD_BRANCH) \
    --build-arg GIT_COMMIT=$(GO_SERVICE_BUILD_COMMIT) \
		.


#
# Docker registry targets
#
docker_registry_push:
	@echo ">>> Tagging docker image with tag."
	@if [ "$(GIT_TAG)" != "" ]; then \
		docker tag $(SERVICE) $(DOCKER_REPO):$(GIT_TAG); \
		docker push $(DOCKER_REPO):$(GIT_TAG); \
	elif [ "$(GIT_BRANCH)" = "master" ]; then \
		docker tag $(SERVICE) $(DOCKER_REPO):latest; \
		docker push $(DOCKER_REPO):latest; \
	else \
		docker tag $(SERVICE) $(DOCKER_REPO):$(GIT_BRANCH); \
		docker push $(DOCKER_REPO):$(GIT_BRANCH); \
	fi