# Variables
BINARY_NAME=gocloak

build:
	mkdir -p bin/$(BINARY_NAME)
	go build -o bin/$(BINARY_NAME) ./cmd/gocloak

run:
	go run ./cmd/gocloak