SHELL := /bin/bash

GOCMD=go
GOMOD=$(GOCMD) mod
GOBUILD=$(GOCMD) build
GOLINT=${GOPATH}/bin/golangci-lint
GORELEASER=/usr/local/bin/goreleaser
GOIMPI=${GOPATH}/bin/impi
GOTEST=$(GOCMD) test

all:
	$(info  "completed running make file for environment-testing garbage collector job")
fmt:
	@go fmt ./...
lint:
	./scripts/lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	@go install github.com/golang/mock/mockgen@latest
	@go install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	@go generate ./...
	$(GOTEST) -short ./... -coverprofile cp.out
integration_test:
	@go get github.com/golang/mock/mockgen@latest
	@go install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	@go generate ./...
	$(GOTEST) ./internal/dao -run TestGarbageCollector_DeleteObsoleteRecordsOneWorkload
	$(GOTEST) ./internal/dao -run TestGarbageCollector_DeleteObsoleteRecordsTwoWorkload
	$(GOTEST) ./internal/dao -run TestGarbageCollector_DeleteEmptyClusters
build:
	$(GOBUILD) -v
build_docker:
	docker build -t gcr.io/snyk-main/sql-raw-data-example:${CIRCLE_SHA1} .
	docker push gcr.io/snyk-main/sql-raw-data-example:${CIRCLE_SHA1}

.PHONY: install-req fmt lint tidy test integration_test build build_docker imports
