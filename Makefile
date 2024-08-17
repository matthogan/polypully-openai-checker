# Define the protobuf file and generated files
PROTO_FILE = api.proto
PROTO_PATH = ./protobuf
GO_OUT = ..

# Define the Go files and packages
SERVICE_FILE = main.go
SERVICE_BINARY = software_info_service
TEST_DIR = ./tests

# Define environment variables for the OpenAI API key
export OPENAI_API_KEY = your_openai_api_key

.PHONY: all proto build run test clean

# Default target
all: proto build test

# Compile the protobuf file
proto:
	@echo "Compiling protobuf files..."
	protoc --go_out=$(GO_OUT) --go-grpc_out=$(GO_OUT) $(PROTO_PATH)/$(PROTO_FILE)

# Build the Go service
# prefer goreleaser...
build: proto
	@echo "Building the service..."
	rm -rf ./dist
	go build -o ./dist/$(SERVICE_BINARY) $(SERVICE_FILE)

# Run the Go service
run: build
	@echo "Running the service..."
	./dist/$(SERVICE_BINARY)

# Run the tests
test:
	@echo "Running tests..."
	go test $(TEST_DIR)

# Clean up generated files and binaries
clean:
	@echo "Cleaning up..."
	rm -rf $(GO_OUT) $(SERVICE_BINARY)

# Re-generate protobuf files, rebuild the service, and run tests
rebuild: clean all

# Check for syntax errors and format Go code
lint:
	@echo "Linting and formatting code..."
	go fmt ./...
	go vet ./...

# Docker related commands
docker-build:
	@echo "Building Docker image..."
	docker build -t software-info-service .

docker-run:
	@echo "Running Docker container..."
	docker run -e OPENAI_API_KEY=$(OPENAI_API_KEY) -p 50051:50051 software-info-service

docker-clean:
	@echo "Cleaning Docker images and containers..."
	docker rmi -f software-info-service
	docker container prune -f

# Help command to list available targets
help:
	@echo "Makefile commands:"
	@echo "  all         - Compile protobuf files, build the service, and run tests"
	@echo "  proto       - Compile protobuf files"
	@echo "  build       - Build the service"
	@echo "  run         - Run the service"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean up generated files and binaries"
	@echo "  rebuild     - Clean, compile, build, and test"
	@echo "  lint        - Check for syntax errors and format Go code"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo "  docker-clean- Clean Docker images and containers"
