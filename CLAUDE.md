# Project Rules

> **Full documentation:** See [AGENTS.md](./AGENTS.md) for complete project guide.

## API Definition Changes

When modifying HTTP API definitions:

1. **Edit `.api` files** in the corresponding project's `desc/` directory:
   - `admin/desc/*.api` for admin service
   - `shop/desc/*.api` for shop service

2. **Run `make api`** to regenerate code:
   ```bash
   cd admin && make api
   # or
   cd shop && make api
   ```

3. **DO NOT edit** auto-generated files:
   - `internal/types/types.go`
   - `internal/handler/routes.go`
   - `internal/handler/*.go`

## Middleware Configuration

**Define middleware in `.api` files, NOT by editing `routes.go`:**

```go
@server (
    group:      users
    middleware: AuthMiddleware  // Add middleware here
)
service admin-api {
    @handler ListUsersHandler
    get /api/v1/users (ListUsersRequest) returns (ListUsersResponse)
}
```

Then implement in `internal/middleware/auth_middleware.go`.

## Build Commands

**ALWAYS use `make build` for compilation:**

```bash
# Build specific service
cd admin && make build
cd shop && make build

# Build all services from root
make build
```

**DO NOT use `go build` directly.**

## After Code Changes

**MUST run `make build`** after any code modifications to verify compilation succeeds.

## Error Handling

**Use custom errors from `pkg/code`, NOT `errors.New()`:**

```go
// BAD: Using errors.New()
var ErrOrderNotFound = errors.New("order not found")
return errors.New("something went wrong")

// GOOD: Use pkg/code errors
return code.ErrOrderNotFound
```

1. **Define all business errors** in `pkg/code/code.go` with proper HTTP status codes and error codes
2. **Error code ranges** by module (see `pkg/code/code.go` for full list):
   - Admin User: 10xxx
   - User: 11xxx
   - Product: 30xxx
   - Order: 40xxx
   - Payment: 50xxx
   - Cart: 60xxx
   - Coupon: 70xxx
   - Promotion: 80xxx
   - Tenant: 90xxx
   - Role: 100xxx
   - Shop: 110xxx
   - Fulfillment: 120xxx
3. **DO NOT create local error variables** with `errors.New()` in application or domain layers