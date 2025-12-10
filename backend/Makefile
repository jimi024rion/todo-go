.PHONY: build run test swagger

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
