# ShopJoy Build Configuration
PROJECT_NAME := shopjoy

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOVET := $(GOCMD) vet

# Directories
ADMIN_DIR := ./admin
SHOP_DIR := ./shop
FRONTEND_DIR := ./shop-admin

# Binaries
ADMIN_BINARY := bin/admin
SHOP_BINARY := bin/shop

.PHONY: all build clean api lint fmt vet test security deps help

# Default target
all: deps build

## build: Build all services
build:
	@echo "Building admin service..."
	cd $(ADMIN_DIR) && $(GOBUILD) -o $(ADMIN_BINARY) .
	@echo "Building shop service..."
	cd $(SHOP_DIR) && $(GOBUILD) -o $(SHOP_BINARY) .

## clean: Clean build artifacts
clean:
	cd $(ADMIN_DIR) && $(GOCLEAN) && rm -rf bin/
	cd $(SHOP_DIR) && $(GOCLEAN) && rm -rf bin/

## api: Generate API code from .api definitions
api:
	cd $(ADMIN_DIR) && make api
	cd $(SHOP_DIR) && make api

## lint: Run golangci-lint
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run --timeout=10m ./...

## fmt: Format Go code
fmt:
	$(GOFMT) -w -s .
	@echo "Go code formatted"

## vet: Run go vet
vet:
	$(GOVET) ./admin/... ./shop/... ./pkg/...

## test: Run all tests
test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

## coverage: Generate test coverage report
coverage: test
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## security: Run security scanner
security:
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securego/gosec/v2/cmd/gosec@latest)
	gosec ./...

## deps: Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) verify

## tidy: Tidy go modules
tidy:
	$(GOMOD) tidy

## frontend-install: Install frontend dependencies
frontend-install:
	cd $(FRONTEND_DIR) && npm install

## frontend-build: Build frontend
frontend-build:
	cd $(FRONTEND_DIR) && npm run build

## frontend-lint: Lint frontend code
frontend-lint:
	cd $(FRONTEND_DIR) && npm run lint

## dev-admin: Run admin service in development
dev-admin:
	cd $(ADMIN_DIR) && go run . -f etc/admin-api.yaml

## dev-shop: Run shop service in development
dev-shop:
	cd $(SHOP_DIR) && go run . -f etc/shop-api.yaml

## dev-frontend: Run frontend development server
dev-frontend:
	cd $(FRONTEND_DIR) && npm run dev

## docker-up: Start Docker services
docker-up:
	docker-compose up -d

## docker-down: Stop Docker services
docker-down:
	docker-compose down

## help: Show this help message
help:
	@echo "ShopJoy - Available targets:"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'