.PHONY: build run test clean deps docker-build deploy

# Build binary
build:
	go build -o bin/paperwork-service cmd/main.go

# Run locally
run:
	go run cmd/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod tidy
	go mod download

# Build Docker image
docker-build:
	docker build -t paperwork-service .

# Deploy to DigitalOcean
deploy:
	doctl apps create-deployment $(shell doctl apps list --format ID --no-header | head -1)

# Development server with auto-reload
dev:
	air -c .air.toml

# Lint code
lint:
	golangci-lint run