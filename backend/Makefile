# ==============================================================================
# VARIABLES
# ==============================================================================

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGENERATE=$(GOCMD) generate
GOMOD=$(GOCMD) mod

# Project variables
BINARY_NAME=./bin/server
CMD_PATH=./cmd/server
# NOTE: `vendor` ディレクトリをテスト対象から除外
PKG_LIST=$(shell go list ./... | grep -v /vendor/)

# Tools
LINTER=golangci-lint
AIR=air
WIRE=wire
BOBGEN=bobgen-psql
ATLAS=atlas

# Help
# INFO: awk スクリプトでヘルプメッセージを生成
# `##@` で始まる行をカテゴリヘッダとして表示
# 参考: https://zenn.dev/loglass/articles/0016-make-makefile
.DEFAULT_GOAL := help
help:
	@awk 'BEGIN {FS = ":.*?##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# ==============================================================================
# PHONY
# ==============================================================================
.PHONY: all help build run test deps clean lint air gen-all gen-wire gen-kessoku gen-bob up down migrate-diff migrate-apply migrate-hash schema-clean swagger aqua-init aqua-generate aqua-install

# ==============================================================================
# MAIN COMMANDS
# ==============================================================================
##@ Application
all: build ## Build the application

build: ## Build the Go application binary
	@echo "Building application..."
	@$(GOBUILD) -o $(BINARY_NAME) $(CMD_PATH)

run: ## Run the application
	@echo "Running application..."
	@$(GORUN) $(CMD_PATH)

test: ## Run all tests
	@echo "Running tests..."
	@$(GOTEST) -v $(PKG_LIST)

deps: ## Install and tidy dependencies
	@echo "Tidying dependencies..."
	@$(GOMOD) tidy

clean: ## Remove build artifacts
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

# ==============================================================================
# DEVELOPMENT COMMANDS
# ==============================================================================
##@ Development
lint: ## Run linter against the codebase
	@echo "Running linter..."
	@$(LINTER) run -v ./...

air: ## Start live-reloading with Air
	@echo "Starting Air for live reloading..."
	@$(AIR) --build.cmd "go build -ldflags '-s -w' -trimpath -o $(BINARY_NAME) $(CMD_PATH)" --build.bin "$(BINARY_NAME)"

# ==============================================================================
# CODE GENERATION COMMANDS
# ==============================================================================
##@ Code Generation
gen-all: gen-wire gen-kessoku gen-bob ## Run all code generators
	@echo "All code generation completed."

gen-wire: ## Generate Wire DI code
	@echo "Generating Wire code..."
	@$(WIRE) ./cmd/server/di/wire

# NOTE: kessoku-migrate と kessoku(go generate) を統合
gen-kessoku: ## Migrate from Wire and generate Kessoku code
	@echo "Migrating and generating Kessoku code..."
# 	@go tool kessoku migrate ./cmd/server/di/wire -o ./cmd/server/di/kessoku/kessoku.go
	@$(GOGENERATE) ./...

gen-bob: ## Generate Bob ORM code
	@echo "Generating Bob ORM code..."
	@$(BOBGEN) generate -c ./bobgen.yaml

gen-swagger: ## Generate Swagger API documentation
	@echo "Generating Swagger docs..."
	@go tool swag init \
		--generalInfo cmd/server/main.go \
		--output docs/swagger \
		--outputTypes go,json,yaml \
		--parseDependency \
		--parseInternal

# ==============================================================================
# DATABASE COMMANDS
# ==============================================================================
##@ Database
up: ## Start database service with Docker Compose
	@echo "Starting Docker Compose services..."
	@docker-compose up -d --build

down: ## Stop database service
	@echo "Stopping Docker Compose services..."
	@docker-compose down

# Use: make migrate-diff name=<migration_name>
migrate-diff: ## Generate a new migration file (e.g., make migrate-diff name=add_users_table)
	@if [ -z "$(name)" ]; then \
		echo "Error: 'name' is a required argument."; \
		echo "Usage: make migrate-diff name=<migration_name>"; \
		exit 1; \
	fi
	@echo "Generating migration file: $(name)..."
	@cd ./internal/infrastructure/rdb && $(ATLAS) migrate diff $(name) --env local

migrate-apply: ## Apply all pending database migrations
	@echo "Applying database migrations..."
	@cd ./internal/infrastructure/rdb && $(ATLAS) migrate apply --env local

migrate-hash: ## Update the migration hash
	@echo "Updating migration hash..."
	@cd ./internal/infrastructure/rdb && $(ATLAS) migrate hash

schema-clean: ## Clean the database schema (destructive)
	@echo "Cleaning database schema..."
	@cd ./internal/infrastructure/rdb && $(ATLAS) schema clean --env local --auto-approve

# ==============================================================================
# TOOLS COMMANDS
# ==============================================================================
##@ Tools
aqua-init: ## Initialize Aqua configuration
	@echo "Initializing Aqua configuration..."
	@aqua init

aqua-generate: ## Interactively generate Aqua configuration
	@echo "Generating Aqua configuration..."
	@aqua g -i -s

aqua-install: ## Install all tools defined in aqua.yaml
	@echo "Installing Aqua dependencies..."
	@aqua i -l
