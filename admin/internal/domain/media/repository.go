package media

import (
	"context"
)

type Repository interface {
	// Insert 插入新资产记录；遇到 uk_provider_public 冲突时返回 ErrDuplicateAsset。
	Insert(ctx context.Context, a *Asset) error

	// FindByID 按主键查询；找不到返回 ErrMediaAssetNotFound（业务错误）。
	FindByID(ctx context.Context, id int64) (*Asset, error)

	// FindByPublicID 按 (provider, public_id) 查。
	FindByPublicID(ctx context.Context, provider, publicID string) (*Asset, error)

	// SoftDelete 软删（GORM 的 deleted_at 由实现处理）。
	SoftDelete(ctx context.Context, id int64) error

	// DeleteByTenant 软删的同时按 tenant_id 过滤。删除成功返回 nil；行不存在
	// 或属于其他租户时统一返回 ErrMediaAssetNotFound，避免 IDOR 信息泄露。
	DeleteByTenant(ctx context.Context, id int64, tenantID int64) error
}