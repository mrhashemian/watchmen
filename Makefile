BASH_PATH:=$(shell which bash)
SHELL=$(BASH_PATH)
ROOT := $(shell realpath $(dir $(lastword $(MAKEFILE_LIST))))
APP := watchmen
BUILD_PATH ?= ".build/app"
BUILD_DATE ?= $(shell TZ="Asia/Tehran" date +'%Y-%m-%dT%H:%M:%S%z')
GIT_BRANCH ?= $(shell echo $(GIT_HEAD_REF) | cut -d'/' -f3)
LDFLAGS := "-w -s"
DC_FILE="docker-compose.yml"
DC_RESOURCE_DIR=".compose"
CURRENT_TIMESTAMP := $(shell date +%s)

all: format lint build-static-vendor

set-goproxy:
	go env -w GOPROXY=direct

############################################################
# Build and Run
############################################################
build:set-goproxy
	go build -v -race .

build-static: set-goproxy
	CGO_ENABLED=0 go build -v -o $(APP) -a -installsuffix cgo -ldflags $(LDFLAGS) .

build-static-vendor: set-goproxy vendor
	CGO_ENABLED=0 go build -mod vendor -v -o $(APP) -installsuffix cgo -ldflags $(LDFLAGS) .

bundle: build-static-vendor
	mkdir -p ${BUILD_PATH}
	cp $(APP) ${BUILD_PATH}

docker: bundle
	DOCKER_BUILDKIT=1 docker build \
	  --build-arg=BUILD_PATH=$(BUILD_PATH) \
	  --build-arg=BUILD_DATE=$(BUILD_DATE) \
	  -t $(APP):latest .

############################################################
# Format and Lint
############################################################
check-goimport:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

format: check-goimport
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s -w R

check-golint:
	which golint || (GO111MODULE=off go get -u -v go get -u golang.org/x/lint/golint)

lint: check-golint
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R golint -set_exit_status R


############################################################
# Development Environment
############################################################
prepare-compose:
	test -d $(DC_RESOURCE_DIR) || mkdir $(DC_RESOURCE_DIR) || true
	test -f $(DC_RESOURCE_DIR)/config.yml || cp config.yml $(DC_RESOURCE_DIR)/config.yml || true

up: prepare-compose docker
	docker-compose up -d

down:
	docker-compose down
