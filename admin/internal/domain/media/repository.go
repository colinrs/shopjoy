package media

import (
	"context"
)

type Repository interface {
	// Insert 插入新资产记录；遇到 uk_provider_public 冲突时返回 ErrDuplicateAsset。
	Insert(ctx context.Context, a *Asset) error

	// FindByID 按主键查询；找不到返回 ErrAssetNotFound（业务错误）。
	FindByID(ctx context.Context, id int64) (*Asset, error)

	// FindByPublicID 按 (provider, public_id) 查。
	FindByPublicID(ctx context.Context, provider, publicID string) (*Asset, error)

	// SoftDelete 软删（GORM 的 deleted_at 由实现处理）。
	SoftDelete(ctx context.Context, id int64) error
}