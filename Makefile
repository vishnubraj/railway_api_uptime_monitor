# Makefile for Railway API Uptime Monitor

.PHONY: build clean run test docker help

# Build configuration
BINARY_NAME=uptime-monitor
CRON_BINARY_NAME=cron-job
BUILD_DIR=build

# Go configuration
GO_FLAGS=-ldflags="-s -w"
GO_CMD=GO111MODULE=on go

help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO_CMD) build $(GO_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Building $(CRON_BINARY_NAME)..."
	@$(GO_CMD) build $(GO_FLAGS) -o $(BUILD_DIR)/$(CRON_BINARY_NAME) ./cmd/cron
	@echo "Build complete!"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME) $(CRON_BINARY_NAME)
	@echo "Clean complete!"

run: ## Run the application
	@echo "Starting uptime monitor..."
	@$(GO_CMD) run main.go

test: ## Run tests
	@echo "Running tests..."
	@$(GO_CMD) test -v ./...

deps: ## Install dependencies
	@echo "Installing dependencies..."
	@$(GO_CMD) mod download
	@$(GO_CMD) mod tidy

fmt: ## Format code
	@echo "Formatting code..."
	@$(GO_CMD) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@$(GO_CMD) vet ./...

docker: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t railway-uptime-monitor .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -d \
		-p 8080:8080 \
		--env-file .env \
		railway-uptime-monitor

dev: ## Run in development mode with auto-reload
	@echo "Starting development server..."
	@which air > /dev/null || $(GO_CMD) install github.com/cosmtrek/air@latest
	@air

setup: ## Setup development environment
	@echo "Setting up development environment..."
	@cp .env.example .env
	@echo "Created .env file - please edit it with your configuration"
	@$(GO_CMD) mod download
	@$(GO_CMD) mod tidy
	@echo "Setup complete!"

deploy: ## Deploy to Railway
	@echo "Deploying to Railway..."
	@./deploy.sh

check: fmt vet ## Run code checks
	@echo "Code checks complete!"

all: clean deps build test ## Build everything

.DEFAULT_GOAL := help
