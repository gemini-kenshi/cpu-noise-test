.PHONY: build clean

BINARY_NAME=cpu-noise-test
BIN_DIR=bin
GOOS=linux
GOARCH=amd64

build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BIN_DIR)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)"

clean:
	@echo "Cleaning $(BIN_DIR) directory..."
	@rm -rf $(BIN_DIR)
	@echo "Clean complete"
