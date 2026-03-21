# Developer Onboarding Guide

Welcome to ShopJoy! This guide will help you set up your development environment and understand the project architecture.

## Prerequisites

- **Go** 1.21+
- **Node.js** 18+ (for frontend development)
- **Docker** & Docker Compose
- **Make** (build automation)
- **MySQL** 8.0+ (or use Docker)
- **Redis** 7+ (or use Docker)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/colinrs/shopjoy.git
cd shopjoy
```

### 2. Start Infrastructure Services

```bash
# Start MySQL and Redis using Docker
docker-compose up -d mysql redis

# Wait for services to be ready (about 10 seconds)
```

### 3. Initialize Database

```bash
# Connect to MySQL and run init script
mysql -h 127.0.0.1 -P 3306 -u root -p123456 < sql/init.sql
```

### 4. Configure Services

Create configuration files based on templates:

```bash
# Admin API config
cp admin/etc/admin-api.yaml.example admin/etc/admin-api.yaml

# Shop API config
cp shop/etc/shop-api.yaml.example shop/etc/shop-api.yaml
```

Edit the YAML files with your local settings.

### 5. Start Backend Services

```bash
# Terminal 1: Admin API
cd admin && go run admin.go -f etc/admin-api.yaml

# Terminal 2: Shop API
cd shop && go run shop.go -f etc/shop-api.yaml
```

### 6. Start Frontend

```bash
# Terminal 3: Admin Dashboard
cd shop-admin
npm install
npm run dev
```

Access the admin dashboard at: http://localhost:5173

## Project Structure

```
shopjoy/
├── admin/                    # Admin backend service
│   ├── desc/                 # API definitions (.api files)
│   ├── internal/
│   │   ├── application/      # Application layer (use cases)
│   │   ├── domain/           # Domain layer (entities, value objects)
│   │   ├── infrastructure/   # Infrastructure layer (repositories)
│   │   ├── handler/          # HTTP handlers (auto-generated)
│   │   ├── logic/            # Business logic handlers
│   │   └── middleware/       # HTTP middleware
│   └── etc/                  # Configuration files
├── shop/                     # Shop frontend API service
├── shop-admin/               # Vue 3 admin dashboard
├── pkg/                      # Shared packages
│   ├── auth/                 # JWT authentication
│   ├── cache/                # Redis caching
│   ├── code/                 # Error codes
│   ├── infra/                # Infrastructure utilities
│   └── snowflake/            # ID generation
└── docs/                     # Documentation
```

## Architecture

ShopJoy follows **Domain-Driven Design (DDD)** with a clean architecture:

```
┌─────────────────────────────────────────┐
│           Interface Layer               │  ← Handlers, Routes
│  (HTTP API, Request/Response types)     │
├─────────────────────────────────────────┤
│         Application Layer               │  ← Services, DTOs
│  (Use case orchestration, Transactions) │
├─────────────────────────────────────────┤
│           Domain Layer                  │  ← Entities, Value Objects
│  (Core business logic, Repository IF)   │     Repository Interfaces
├─────────────────────────────────────────┤
│       Infrastructure Layer              │  ← Repository Implementations
│  (Database, External Services, Cache)   │     External API clients
└─────────────────────────────────────────┘
```

### Key Patterns

1. **Repository Pattern**: Data access abstraction
2. **Multi-tenancy**: Tenant isolation via `tenant_id`
3. **Value Objects**: `Money`, `TenantID`, `AuditInfo`
4. **Domain Events**: For cross-context communication

## Development Workflow

### Modifying API Definitions

1. Edit `.api` files in `admin/desc/` or `shop/desc/`
2. Run code generation:
   ```bash
   cd admin && make api
   # or
   cd shop && make api
   ```
3. Implement the generated logic in `internal/logic/`

### Adding a New Feature

1. **Define the domain entity** in `internal/domain/`
2. **Create repository interface** in the domain package
3. **Implement the repository** in `internal/infrastructure/persistence/`
4. **Create application service** if complex orchestration needed
5. **Define API** in `desc/*.api` and run `make api`
6. **Implement logic handler** in `internal/logic/`

### Code Standards

- **Go**: Follow [Effective Go](https://golang.org/doc/effective_go) and [Uber Style Guide](https://github.com/uber-go/guide)
- **TypeScript/Vue**: Follow [Vue Style Guide](https://vuejs.org/style-guide/)
- **Commit Messages**: Use conventional commits format

```bash
# Run linters
cd admin && make lint  # (if configured)
cd shop-admin && npm run lint
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./admin/internal/domain/product/...
```

## Common Tasks

### Create a New Admin User

```bash
# Using the API
curl -X POST http://localhost:8888/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "real_name": "Admin User",
    "password": "password123"
  }'
```

### View API Documentation

API definitions are in the `desc/` directories. Each `.api` file documents the endpoints.

## Troubleshooting

### Database Connection Failed

```bash
# Check MySQL is running
docker-compose ps mysql

# Check connection
mysql -h 127.0.0.1 -P 3306 -u root -p123456
```

### Redis Connection Failed

```bash
# Check Redis is running
docker-compose ps redis

# Test connection
redis-cli -h 127.0.0.1 -p 6379 ping
```

### Port Already in Use

```bash
# Find process using port 8888
lsof -i :8888

# Kill the process
kill -9 <PID>
```

## Useful Commands

```bash
# Build all services
make build

# Build specific service
cd admin && make build

# Generate API code
cd admin && make api

# Run with hot reload (requires air)
air -c .air.toml

# Check for race conditions
go test -race ./...
```

## Getting Help

- **Architecture**: See `docs/ARCHITECTURE.md`
- **API Reference**: See `docs/API.md`
- **Issues**: Create an issue in the repository
- **Code Questions**: Check inline comments and godoc

## Next Steps

1. Read the [Architecture Documentation](./ARCHITECTURE.md)
2. Review the [API Documentation](./API.md)
3. Explore the domain models in `admin/internal/domain/`
4. Try adding a simple feature following the workflow above