.PHONY: build run test swagger up down migrate-up

# Go configuration
BINARY_NAME=todo-go
# Path is relative to the 'backend' directory
CMD_PATH=./cmd/server

# Build the application
build:
	@echo "Building Go binary..."
	go build -o ../$(BINARY_NAME) $(CMD_PATH)

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
	air --build.cmd "go build -ldflags '-s -w' -trimpath -o ./bin/server ./cmd/server/main.go" --build.bin "./bin/server"

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

# Air live reloadr
air-reload:
	air --build.cmd "go build -ldflags '-s -w' -trimpath -o ./bin/server ./cmd/server/main.go" --build.bin "./bin/server"

# wire generate
wire:
	wire ./cmd/server/di
	@echo "Wire code generation completed."
