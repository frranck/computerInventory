BINARY_NAME=computerInventory
PACKAGE_PATH=./cmd/main.go

.PHONY: build run test docker docker-up docker-down clean

## Build the binary
build:
	go build -o $(BINARY_NAME) $(PACKAGE_PATH)

## Run locally
run: build
	./$(BINARY_NAME)

## Run all tests
test:
	go test ./... -v

## Build Docker image
docker:
	docker build -t $(BINARY_NAME) .

## Run via Docker Compose
docker-up:
	docker compose up --build

## Stop and remove containers
docker-down:
	docker compose down

## Clean up
clean:
	rm -f $(BINARY_NAME)
