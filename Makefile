PROJECT_NAME := ipmi_power_exporter
BUILD_DIR := build

# Binary targets
BINARY_NAME := $(BUILD_DIR)/$(PROJECT_NAME)

# Go build flags
GO_BUILD_FLAGS := -o $(BINARY_NAME)

# Go test flags
GO_TEST_FLAGS := -v -race ./...

# Go clean flags
GO_CLEAN_FLAGS := -i

# Go install flags
GO_INSTALL_FLAGS :=

# Go get flags
GO_GET_FLAGS := -v

# Go list flags
GO_LIST_FLAGS := -f '{{join .Deps "\n"}}'

# Go list... targets
GO_LIST_TARGETS := $(shell go list $(GO_LIST_FLAGS))

# Go mod tidy flags
GO_MOD_TIDY_FLAGS := -v

# Go mod vendor flags
GO_MOD_VENDOR_FLAGS :=

# Go mod verify flags
GO_MOD_VERIFY_FLAGS :=

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
	if (helpMessage) { \
		helpCommand = substr($$1, 0, index($$1, ":")-1); \
		helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
		printf "\033[36m%-20s\033[0m %s\n", helpCommand, helpMessage; \
	} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

# Build target
.PHONY: build
build:
	@echo "Building $(PROJECT_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(GO_BUILD_FLAGS)

# Test target
.PHONY: test
test:
	@echo "Running tests..."
	@go test $(GO_TEST_FLAGS)

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@go clean $(GO_CLEAN_FLAGS)
	@rm -rf $(BUILD_DIR)

# Install target
.PHONY: install
install:
	@echo "Installing $(PROJECT_NAME)..."
	@go install $(GO_INSTALL_FLAGS)

# Get target
.PHONY: get
get:
	@echo "Getting dependencies..."
	@go get $(GO_GET_FLAGS) $(GO_LIST_TARGETS)

# Mod tidy target
.PHONY: mod-tidy
mod-tidy:
	@echo "Tidying go.mod and go.sum..."
	@go mod tidy $(GO_MOD_TIDY_FLAGS)

# Mod vendor target
.PHONY: mod-vendor
mod-vendor:
	@echo "Vendoring dependencies..."
	@go mod vendor $(GO_MOD_VENDOR_FLAGS)

# Mod verify target
.PHONY: mod-verify
mod-verify:
	@echo "Verifying dependencies..."
	@go mod verify $(GO_MOD_VERIFY_FLAGS)

# All target
.PHONY: all
all: build test
