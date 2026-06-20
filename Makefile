.PHONY: help build run test clean lint fmt docker-build docker-run dev

help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  run           - Run the server"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  lint          - Run linters"
	@echo "  fmt           - Format code"
	@echo "  dev           - Run in development mode"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker-push   - Push Docker image"

build:
	go build -o mcp-tududi ./cmd/server

run: build
	./mcp-tududi

test:
	go test -v -race -coverprofile=coverage.out ./...

clean:
	go clean
	rm -f mcp-tududi coverage.out

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...
	goimports -w .

dev:
	go run ./cmd/server

docker-build:
	docker build -t mcp-tududi:latest .

docker-run: docker-build
	docker run -it --env-file .env -p 8080:8080 mcp-tududi:latest

docker-push:
	docker push ghcr.io/nouchka/mcp-tududi:latest

deps:
	go mod download
	go mod verify
