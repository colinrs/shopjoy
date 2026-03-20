package contextx

import (
	"context"
)

type contextKey string

const (
	userIDKey    contextKey = "user_id"
	tenantIDKey  contextKey = "tenant_id"
	userTypeKey  contextKey = "user_type"
)

// SetUserID 设置用户ID到 context
func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID 从 context 获取用户ID
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}

// SetTenantID 设置租户ID到 context
func SetTenantID(ctx context.Context, tenantID int64) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// GetTenantID 从 context 获取租户ID
func GetTenantID(ctx context.Context) (int64, bool) {
	tenantID, ok := ctx.Value(tenantIDKey).(int64)
	return tenantID, ok
}

// SetUserType 设置用户类型到 context
func SetUserType(ctx context.Context, userType int) context.Context {
	return context.WithValue(ctx, userTypeKey, userType)
}

// GetUserType 从 context 获取用户类型
func GetUserType(ctx context.Context) (int, bool) {
	userType, ok := ctx.Value(userTypeKey).(int)
	return userType, ok
}

// GetCurrentUserID 获取当前用户ID，如果不存在返回0
func GetCurrentUserID(ctx context.Context) int64 {
	userID, _ := GetUserID(ctx)
	return userID
}