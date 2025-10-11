.PHONY: run build test clean dev docker-up docker-down

# Run the application
run:
	go run main.go

# Build the application
build:
	go build -o lacpa

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f lacpa
	go clean

# Run with hot reload (requires air)
dev:
	air

# Start MongoDB with Docker Compose
docker-up:
	docker-compose up -d

# Stop MongoDB
docker-down:
	docker-compose down

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Download dependencies
deps:
	go mod download

# Tidy dependencies
tidy:
	go mod tidy

# Install development tools
install-tools:
	go install github.com/air-verse/air@latest

# Full setup (dependencies + docker)
setup: deps docker-up
	@echo "Setup complete! Run 'make run' to start the server"
