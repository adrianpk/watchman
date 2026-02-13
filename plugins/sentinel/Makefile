BIN_NAME = sentinel

all: build

help:
	@echo "Available targets:"
	@echo ""
	@echo "  build     - Build the binary"
	@echo "  install   - Install to GOPATH/bin"
	@echo "  test      - Run tests"
	@echo "  clean     - Remove binary"
	@echo "  tidy      - Run go mod tidy"

build:
	@echo "Building $(BIN_NAME)..."
	@go build -o $(BIN_NAME) .

install:
	@echo "Installing $(BIN_NAME)..."
	@go install .

test:
	@go test ./...

clean:
	@rm -f $(BIN_NAME)

tidy:
	@go mod tidy

.PHONY: all build install test clean tidy help
