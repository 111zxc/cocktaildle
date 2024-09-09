BINARY_NAME=cocktaildle_backend
SRC=$(shell find . -type f -name '*.go')

all: fmt vet build

fmt:
	@echo "Running go fmt..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

build: fmt vet
	@echo "Building the application..."
	@go build -o $(BINARY_NAME) cmd/main.go

run: build
	@echo "Running the application..."
	@./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test ./...

.PHONY: all fmt vet build run clean test
