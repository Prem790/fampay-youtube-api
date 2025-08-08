.PHONY: build run test clean docker-up docker-down deps fmt

# Build the application
build:
	go build -o bin/main cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with MongoDB
test-integration:
	# Make sure MongoDB is running on localhost:27017
	go test -v ./internal/repository/
	go test -v ./internal/api/handlers/

# Clean build artifacts
clean:
	rm -rf bin/

# Docker commands for MongoDB
docker-up:
	docker-compose -f docker-compose.mongodb.yml up -d

docker-down:
	docker-compose -f docker-compose.mongodb.yml down

docker-build:
	docker-compose -f docker-compose.mongodb.yml build

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# MongoDB specific commands
mongo-shell:
	docker exec -it $(docker ps -q -f name=mongodb) mongosh fampay_youtube

mongo-status:
	docker-compose -f docker-compose.mongodb.yml ps

# Load test data into MongoDB
load-test-data:
	@echo "Loading test data into MongoDB..."
	@go run scripts/load_test_data.go