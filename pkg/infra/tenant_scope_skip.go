package infra

import "context"

type skipTenantScopeKeyType struct{}

var skipTenantScopeKey = skipTenantScopeKeyType{}

// SkipTenantScope 返回一个跳过租户过滤的 context
// 用于平台管理员全量查询、数据迁移等需要绕过租户隔离的场景
func SkipTenantScope(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipTenantScopeKey, true)
}

// IsSkipTenantScope 检查 context 是否标记为跳过租户过滤
func IsSkipTenantScope(ctx context.Context) bool {
	v, _ := ctx.Value(skipTenantScopeKey).(bool)
	return v
}
