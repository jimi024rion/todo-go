.PHONY: build run test swagger clean bobgen

APP_NAME := todo-api
CMD_DIR := ./cmd/app
BUILD_DIR := ./bin
SWAGGER_OUTPUT_DIR := ./docs

build:
	@echo "Building $(APP_NAME)..."
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

run: build
	@echo "Running $(APP_NAME)..."
	$(BUILD_DIR)/$(APP_NAME)

test:
	@echo "Running tests..."
	go test ./...

swagger:
	@echo "Generating Swagger documentation (placeholder)..."
	swag init -g ./cmd/app/main.go -o $(SWAGGER_OUTPUT_DIR)

bobgen:
	@echo "Generating Bob ORM models..."
	@bobgen-psql --config bobgen.yaml

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(SWAGGER_OUTPUT_DIR)/swagger.json
	@rm -f $(SWAGGER_OUTPUT_DIR)/swagger.yaml
	@rm -rf models dbinfo factory dberrors
