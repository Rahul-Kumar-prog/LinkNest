# LinkNest Project Makefile
# A comprehensive Makefile for managing the entire project

# Colors for output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[1;33m
NC = \033[0m # No Color

.PHONY: help install deps-backend deps-frontend run dev build clean start stop status logs

# Default target
all: help

help: ## Show this help message
	@echo ""
	@echo "LinkNest Makefile - Available Commands:"
	@echo ""
	@echo "  $(GREEN)make install$(NC)         - Install all dependencies (backend + frontend)"
	@echo "  $(GREEN)make deps-backend$(NC)    - Install backend dependencies only"
	@echo "  $(GREEN)make deps-frontend$(NC)   - Install frontend dependencies only"
	@echo "  $(GREEN)make dev$(NC)             - Run both frontend and backend in dev mode"
	@echo "  $(GREEN)make run$(NC)             - Run both frontend and backend (production)"
	@echo "  $(GREEN)make build$(NC)           - Build all projects"
	@echo "  $(GREEN)make clean$(NC)           - Clean build artifacts"
	@echo "  $(GREEN)make start$(NC)           - Start all services"
	@echo "  $(GREEN)make stop$(NC)            - Stop all services"
	@echo "  $(GREEN)make status$(NC)          - Show status of all services"
	@echo "  $(GREEN)make logs$(NC)            - Show logs from all services"
	@echo "  $(GREEN)make help$(NC)            - Show this help message"
	@echo ""

install: deps-backend deps-frontend ## Install all dependencies
	@echo "$(GREEN)✅ All dependencies installed successfully!$(NC)"

deps-backend: ## Install backend dependencies
	@echo "$(YELLOW)📦 Installing backend dependencies...$(NC)"
	@cd backend && go mod download
	@echo "$(GREEN)✅ Backend dependencies installed!$(NC)"

deps-frontend: ## Install frontend dependencies
	@echo "$(YELLOW)📦 Installing frontend dependencies...$(NC)"
	@cd frontend && npm install
	@echo "$(GREEN)✅ Frontend dependencies installed!$(NC)"

dev: ## Run both frontend and backend in development mode (using background processes)
	@echo "$(GREEN)🚀 Starting development servers...$(NC)"
	@echo "$(YELLOW)Backend: http://localhost:8080$(NC)"
	@echo "$(YELLOW)Frontend: http://localhost:5173$(NC)"
	@echo ""
	@make -j2 run-backend run-frontend

run-frontend: ## Run frontend dev server
	@echo "$(YELLOW)🎨 Starting frontend dev server...$(NC)"
	@cd frontend && npm run dev

run-backend: ## Run backend server
	@echo "$(YELLOW)🔧 Starting backend server...$(NC)"
	@cd backend && go run cmd/server/main.go

run: ## Run both frontend and backend
	@echo "$(GREEN)🚀 Starting all services...$(NC)"
	@make -j2 run-frontend run-backend

build: ## Build all projects
	@echo "$(GREEN)🔨 Building all projects...$(NC)"
	@make build-backend
	@make build-frontend
	@echo "$(GREEN)✅ All projects built successfully!$(NC)"

build-backend: ## Build backend binary
	@echo "$(YELLOW)🔧 Building backend...$(NC)"
	@cd backend && go build -o bin/server cmd/server/main.go
	@echo "$(GREEN)✅ Backend built successfully!$(NC)"

build-frontend: ## Build frontend for production
	@echo "$(YELLOW)🎨 Building frontend...$(NC)"
	@cd frontend && npm run build
	@echo "$(GREEN)✅ Frontend built successfully!$(NC)"

start: ## Start built applications
	@echo "$(GREEN)🚀 Starting all services...$(NC)"
	@./backend/bin/server &
	@echo "$(YELLOW)Backend started on http://localhost:8080$(NC)"

clean: ## Clean all build artifacts
	@echo "$(YELLOW)🧹 Cleaning build artifacts...$(NC)"
	@rm -rf backend/bin
	@rm -rf frontend/dist
	@rm -rf frontend/node_modules
	@echo "$(GREEN)✅ Cleaned successfully!$(NC)"

