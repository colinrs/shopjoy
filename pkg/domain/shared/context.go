package shared

import "context"

type contextKey string

const (
	tenantIDKey contextKey = "tenant_id"
	userIDKey   contextKey = "user_id"
)

// GetTenantIDFromContext extracts tenant ID from context
func GetTenantIDFromContext(ctx context.Context) TenantID {
	if v := ctx.Value(tenantIDKey); v != nil {
		if tid, ok := v.(int64); ok {
			return TenantID(tid)
		}
	}
	return 0
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) int64 {
	if v := ctx.Value(userIDKey); v != nil {
		if uid, ok := v.(int64); ok {
			return uid
		}
	}
	return 0
}

// SetTenantIDInContext sets tenant ID in context
func SetTenantIDInContext(ctx context.Context, tenantID int64) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// SetUserIDInContext sets user ID in context
func SetUserIDInContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}