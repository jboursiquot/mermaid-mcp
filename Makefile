default: help

help: ## show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

tidy: ## go mod tidy
	go mod tidy

vendor: ## vendor dependencies
	go mod vendor

test: ## run tests
	go test ./... -v -race

build: ## build the project
	go build -mod=vendor -o bin/server ./cmd/server/*.go
	
build-with-vendor: vendor build ## vendor dependencies and build

run: ## run the project
	go run -mod=vendor ./cmd/server/*.go

debug: ## build and run the container in debug mode
	docker build -t mermaid-mcp:debug .
	docker run \
		-e MERMAID_MCP_PG_CONN_STR=${MERMAID_MCP_PG_CONN_STR} \
		--rm \
		-it \
		-p 6274:6274 \
		-p 6277:6277 \
		mermaid-mcp:debug

DATABASE_URL ?= set-me

up: vendor ## start docker containers
	DATABASE_URL=${DATABASE_URL} docker-compose up --build

down: ## stop docker containers
	docker-compose down