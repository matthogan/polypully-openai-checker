# Define the protobuf file and generated files
PROTO_FILE = api.proto
PROTO_PATH = ./protobuf
GO_OUT = ..

# Define the Go files and packages
SERVICE_FILE = main.go
SERVICE_BINARY = software_info_service
TEST_DIR = ./tests

# Define file names for certificates and keys
CA_KEY = .secret/ca.key
CA_CERT = .secret/ca.crt
SERVER_KEY = .secret/server.key
SERVER_CSR = .secret/server.csr
SERVER_CERT = .secret/server.crt
CLIENT_KEY = .secret/client.key
CLIENT_CSR = .secret/client.csr
CLIENT_CERT = .secret/client.crt

# Define environment variables for the OpenAI API key
export OPENAI_API_KEY = your_openai_api_key

.PHONY: all proto build run test clean

# Default target
all: proto build test

# update all dependencies
update:
	go get -u all

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

# Run the tests
test: build
	@echo "Running tests..."
	go test -v ./...

# Run the Go service
run: build test
	@echo "Running the service..."
	./dist/$(SERVICE_BINARY)

# Clean up generated files and binaries
clean:
	@echo "Cleaning up..."
	rm -rf ./dist

# Re-generate protobuf files, rebuild the service, and run tests
rebuild: clean all

# Sanity test a local instance
sanity-test:
	@echo "Sanity testing the service..."
	@echo "Sending a request to the service..."
	@grpcurl -cert .secret/client.crt -key .secret/client.key -cacert .secret/ca.crt -d '{"name":"Notepad++"}' localhost:50051 software.SoftwareInfoService/GetSoftwareInfo
	@echo

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
	@echo "Cleaning docker images and containers..."
	docker rmi -f software-info-service
	docker container prune -f

docker-prometheus:
	@echo "Running prometheus docker container..."
	docker run --name prometheus --rm -v "$(shell pwd)/metrics/prometheus.yml:/etc/prometheus/prometheus.yml" -p 9090:9090 prom/prometheus

docker-memcached:
	@echo "Running memcached docker container..."
	docker run --name memcached --rm -p 11211:11211 memcached

# Generate CA key and self-signed certificate
ca-cert:
	@echo "Generating CA key and self-signed certificate..."
	openssl genpkey -algorithm RSA -out $(CA_KEY) -pkeyopt rsa_keygen_bits:2048
	openssl req -x509 -new -nodes -key $(CA_KEY) -sha256 -days 365 -out $(CA_CERT) -subj "/C=US/ST=State/L=City/O=Organization/OU=OrgUnit/CN=localhost"

# Generate server key, CSR, and certificate signed by CA with SANs
server-cert:
	@echo "Generating server key, CSR, and certificate signed by CA with SANs..."
	openssl genpkey -algorithm RSA -out $(SERVER_KEY) -pkeyopt rsa_keygen_bits:2048
	openssl req -new -key $(SERVER_KEY) -out $(SERVER_CSR) -subj "/C=US/ST=State/L=City/O=Organization/OU=OrgUnit/CN=localhost" -addext "subjectAltName = DNS:localhost"
	echo "subjectAltName=DNS:localhost" > /tmp/server_extfile.cnf
	openssl x509 -req -in $(SERVER_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -out $(SERVER_CERT) -days 365 -sha256 -extfile /tmp/server_extfile.cnf
	rm /tmp/server_extfile.cnf

# Generate client key, CSR, and certificate signed by CA with SANs
client-cert:
	@echo "Generating client key, CSR, and certificate signed by CA with SANs..."
	openssl genpkey -algorithm RSA -out $(CLIENT_KEY) -pkeyopt rsa_keygen_bits:2048
	openssl req -new -key $(CLIENT_KEY) -out $(CLIENT_CSR) -subj "/C=US/ST=State/L=City/O=Organization/OU=OrgUnit/CN=localhost" -addext "subjectAltName = DNS:localhost"
	echo "subjectAltName=DNS:localhost" > /tmp/client_extfile.cnf
	openssl x509 -req -in $(CLIENT_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -out $(CLIENT_CERT) -days 365 -sha256 -extfile /tmp/client_extfile.cnf
	rm /tmp/client_extfile.cnf

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
	@echo "  ca-cert     - Generate CA key and self-signed certificate (dev)"
	@echo "  server-cert - Generate server key, CSR, and certificate signed by CA (dev)"
	@echo "  client-cert - Generate client key, CSR, and certificate signed by CA (dev)"
