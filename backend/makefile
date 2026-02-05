# Backend Makefile
# Manage Go backend development and operations

# Colors for output
GREEN = \033[0;32m
YELLOW = \033[1;33m
NC = \033[0m

.PHONY: help run build test clean deps lint fmt run-air

# Default target
all: help

help: ## Show this help message
	@echo ""
	@echo "Backend Makefile - Available Commands:"
	@echo ""
	@echo "  $(GREEN)make run$(NC)             - Run backend server directly"
	@echo "  $(GREEN)make run-air$(NC)         - Run with hot-reload (using Air)"
	@echo "  $(GREEN)make build$(NC)           - Build backend binary"
	@echo "  $(GREEN)make test$(NC)            - Run all tests"
	@echo "  $(GREEN)make clean$(NC)           - Clean build artifacts"
	@echo "  $(GREEN)make deps$(NC)            - Download dependencies"
	@echo "  $(GREEN)make lint$(NC)            - Run linter"
	@echo "  $(GREEN)make fmt$(NC)             - Format Go code"
	@echo "  $(GREEN)make help$(NC)            - Show this help message"
	@echo ""

run: ## Run backend server directly
	@echo "$(YELLOW)🔧 Starting backend server...$(NC)"
	@go run cmd/server/main.go

run-air: ## Run with hot-reload using Air
	@echo "$(YELLOW)🔥 Starting with hot-reload (Air)...$(NC)"
	@air

build: ## Build backend binary
	@echo "$(YELLOW)🔨 Building backend...$(NC)"
	@go build -o bin/server cmd/server/main.go
	@echo "$(GREEN)✅ Backend binary built at bin/server$(NC)"

test: ## Run all tests
	@echo "$(YELLOW)🧪 Running tests...$(NC)"
	@go test ./...
	@echo "$(GREEN)✅ All tests passed!$(NC)"

clean: ## Clean build artifacts
	@echo "$(YELLOW)🧹 Cleaning...$(NC)"
	@rm -rf bin/
	@go clean
	@echo "$(GREEN)✅ Cleaned successfully!$(NC)"

deps: ## Download dependencies
	@echo "$(YELLOW)📦 Downloading dependencies...$(NC)"
	@go mod download
	@echo "$(GREEN)✅ Dependencies downloaded!$(NC)"

lint: ## Run linter (requires golangci-lint)
	@echo "$(YELLOW)🔍 Running linter...$(NC)"
	@golangci-lint run ./...
	@echo "$(GREEN)✅ Linting passed!$(NC)"

fmt: ## Format Go code
	@echo "$(YELLOW)📝 Formatting code...$(NC)"
	@gofmt -w .
	@goimports -w .
	@echo "$(GREEN)✅ Code formatted!$(NC)"

