package contextx

import (
	"context"
	"net/http"

	"github.com/colinrs/shopjoy/pkg/code"
)

type contextKey string

const (
	userIDKey   contextKey = "user_id"
	tenantIDKey contextKey = "tenant_id"
	userTypeKey contextKey = "user_type"
	userNameKey contextKey = "user_name"
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

// MustGetTenantID 获取租户ID，如果不存在返回错误
// 用于需要确保租户ID存在的场景
func MustGetTenantID(ctx context.Context) (int64, error) {
	tenantID, ok := GetTenantID(ctx)
	if !ok || tenantID == 0 {
		return 0, ErrTenantNotFound
	}
	return tenantID, nil
}

// MustGetUserID 获取用户ID，如果不存在返回错误
func MustGetUserID(ctx context.Context) (int64, error) {
	userID, ok := GetUserID(ctx)
	if !ok || userID == 0 {
		return 0, ErrUserNotFound
	}
	return userID, nil
}

// GetCurrentUserType 获取当前用户类型，如果不存在返回0
func GetCurrentUserType(ctx context.Context) int {
	userType, _ := GetUserType(ctx)
	return userType
}

// SetUserName 设置用户名到 context
func SetUserName(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, userNameKey, userName)
}

// GetUserName 从 context 获取用户名
func GetUserName(ctx context.Context) (string, bool) {
	userName, ok := ctx.Value(userNameKey).(string)
	return userName, ok
}

// GetCurrentUserName 获取当前用户名，如果不存在返回空字符串
func GetCurrentUserName(ctx context.Context) string {
	userName, _ := GetUserName(ctx)
	return userName
}

// IsPlatformAdmin 检查当前用户是否为平台管理员
func IsPlatformAdmin(ctx context.Context) bool {
	userType, ok := GetUserType(ctx)
	return ok && userType == 1 // 1 = 平台超管
}

// Context errors using pkg/code
var (
	ErrTenantNotFound = &code.Err{HTTPCode: http.StatusBadRequest, Code: 90001, Msg: "tenant not found in context"}
	ErrUserNotFound   = &code.Err{HTTPCode: http.StatusBadRequest, Code: 11004, Msg: "user not found in context"}
)