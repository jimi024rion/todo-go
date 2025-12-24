.PHONY: build run test swagger up down migrate-up migrate-diff migrate-apply migrate-hash db-reset schema-clean

# Go configuration
BINARY_NAME=./bin/server
# Path is relative to the 'backend' directory
CMD_PATH=./cmd/server

# Build the application
build:
	@echo "Building Go binary..."
	go build -o $(BINARY_NAME) $(CMD_PATH)

# Run the application
run:
	@echo "Running Go application..."
	go run $(CMD_PATH)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Generate Swagger docs
swagger:
	@echo "Generating Swagger docs..."
	@echo "swag init will be configured here later"

air:
	@echo "Starting Air for live reloading..."
	air --build.cmd "go build -ldflags '-s -w' -trimpath -o $(BINARY_NAME) $(CMD_PATH)" --build.bin "$(BINARY_NAME)"

# Start Docker Compose services
up:
	@echo "Starting Docker Compose services..."
	docker-compose up -d --build

# Stop Docker Compose services
down:
	@echo "Stopping Docker Compose services..."
	docker-compose down

# Run database migrations
migrate-up:
	@echo "Running database migrations..."
	@echo "Atlas migration command will be configured here later"

# wire generate
wire:
	wire ./cmd/server/di
	@echo "Wire code generation completed."

bob:
	bobgen-psql generate -c ./bobgen.yaml
	@echo "Bob code generation completed."

# migrate-diff 用の引数処理: make migrate-diff <name>
ifeq (migrate-diff,$(firstword $(MAKECMDGOALS)))
  MIGRATION_NAME := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(MIGRATION_NAME):;@:)
endif

migrate-diff:
	cd ./internal/infrastructure/rdb && atlas migrate diff $(MIGRATION_NAME) --env local
	@echo "Atlas migration diff completed."

migrate-apply:
	cd ./internal/infrastructure/rdb && atlas migrate apply --env local
	@echo "Atlas migration apply completed."

migrate-hash:
	cd ./internal/infrastructure/rdb && atlas migrate hash
	@echo "Atlas migration hash updated."

schema-clean:
	cd ./internal/infrastructure/rdb && atlas schema clean --env local --auto-approve
	@echo "Atlas schema clean completed."

aqua-init:
	@echo "Initializing Aqua configuration..."
	aqua init

aqua-generate:
	@echo "Generating Aqua configuration..."
	aqua g -i -s

aqua-install:
	@echo "Installing Aqua dependencies..."
	aqua i -l
