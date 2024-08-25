DEFAULT_GOAL := help

BUILD_FOLDER = dist
CRT_FOLDER = ssl/ca

# Build info
CLIENT_VERSION ?= 0.1.0

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: proto
proto: ## Generate gRPC protobuf bindings
	./scripts/gen-proto

.PHONY: keeper ## Build the gophkeeper service
keeper:
	go build -o $(BUILD_FOLDER)/$@ cmd/$@/*.go

.PHONY: keeperctl ## Build the gophkeeper client
keepctl:
	./scripts/build-client $(CLIENT_VERSION)

.PHONY: all ## Build whole product.
all: keeper keeperctl

.PHONY: download
download: ## Download go.mod dependencies
	echo Downloading go.mod dependencies
	go mod download

.PHONY: clean
clean:
	rm -rf $(BUILD_FOLDER) $(CRT_FOLDER)

.PHONY: install-tools
install-tools: ## Install additional linters and dev tools
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest

.PHONY: lint
lint: ## Run linters on the source code
	golangci-lint run

.PHONY: test
test: ## Run unit tests
	@go test -v -race ./... -coverprofile=coverage.out.tmp -covermode atomic
	@cat coverage.out.tmp | grep -v -E "(_mock|.pb).go" > coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out

.PHONY: update-snapshots
update-snapshots: ## Update unit-tests's snapshots
	@UPDATE_SNAPS=true go test -v -race ./... -coverprofile=coverage.out.tmp -covermode atomic

.PHONY: ssl
ssl: ## Generate SSL certificates for secure communications
	./scripts/gen-ca
	./scripts/issue-crt