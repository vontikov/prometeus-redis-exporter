PROJECTNAME = $(shell basename "$(PWD)")

include project.properties

BASE_DIR   = $(shell pwd)
ASSET_DIR  = $(BASE_DIR)/assets
BIN_DIR    = $(BASE_DIR)/.bin
TMP_DIR    = $(BASE_DIR)/.tmp

GO_FILES   = $(shell find . -type f -name '*.go')
GO_LDFLAGS = "-s -w -X main.App=$(APP) -X main.Version=$(VERSION)"

TEST_OPTS  = -timeout 300s -failfast -v

TEST_PKGS  = \
  ./internal/parser

BINARY_FILE   = $(BIN_DIR)/$(APP)
COVERAGE_FILE = $(BIN_DIR)/test-coverage.out

BIN_MARKER   = $(TMP_DIR)/bin.marker
TEST_MARKER  = $(TMP_DIR)/test.marker

ifndef VERBOSE
.SILENT:
endif

.PHONY: all help clean deps

all: help

## help: Prints help
help: Makefile
	@echo "Choose a command in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## test: Runs tests.
test: $(TEST_MARKER)

$(TEST_MARKER): prepare $(GO_FILES)
	@echo "Executing tests..."
	@GOBIN=$(BIN_DIR); go test $(TEST_OPTS) $(TEST_PKGS)
	@mkdir -p $(@D)
	@touch $@

## build: Builds the binaries
build: prepare $(BIN_MARKER)

$(BIN_MARKER): $(GO_FILES)
	@echo "Building the binaries..."
	@GOBIN=$(BIN_DIR) \
      CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
      go build -tags=$(tags) -ldflags=$(GO_LDFLAGS) -o $(BINARY_FILE) cmd/*.go
	@mkdir -p $(@D)
	@touch $@

## clean: Cleans the project
clean:
	@echo "Cleaning the project..."
	@rm -rf $(BIN_DIR) $(TMP_DIR)

## image: Builds the Docker image
image: build
	@echo "Building Docker image..."
	$(eval tmp := $(shell mktemp -d))
	@cp -r $(BINARY_FILE) $(tmp)
	@cp -r $(ASSET_DIR)/docker/Dockerfile $(tmp)
	@docker build --build-arg VERSION=$(VERSION) \
      -t $(DOCKER_NS)/$(APP):$(VERSION) $(tmp)

## deps: Installs dependencies
deps:
	@go get -u github.com/go-redis/redis/v8
	@go get -u github.com/gorilla/mux
	@go get -u github.com/hashicorp/go-hclog
	@go get -u github.com/prometheus/client_golang/prometheus/promhttp
	@go get -u github.com/sirupsen/logrus
	@go get -u github.com/stretchr/testify/assert

prepare:
	@mkdir -p $(BIN_DIR)
	@mkdir -p $(TMP_DIR)
