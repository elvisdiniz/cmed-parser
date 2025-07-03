BINARY_NAME=cmed-parser
DIST_DIR=dist/bin

.PHONY: all build clean test build-linux build-windows build-macos

all: test build

test:
	@echo "Running tests..."
	@go test ./...

build: build-linux build-windows build-macos

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(DIST_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 main.go

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(DIST_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go

build-macos:
	@echo "Building for macOS..."
	@mkdir -p $(DIST_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-macos-amd64 main.go

clean:
	@echo "Cleaning up..."
	@rm -rf dist
