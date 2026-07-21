package region

import (
	"context"

	"gorm.io/gorm"
)

// RegionQuery 区域查询条件
type RegionQuery struct {
	TenantID    int64  // 0 表示只查询平台预置数据
	CountryCode string // ISO 3166-1 alpha-2，留空表示不按国家过滤
	ParentCode  string // 父级区域 code，留空表示查询顶级（Level=1）
	Level       int    // 1=国家, 2=省, 3=市, 4=区；0 表示不按层级过滤
	IsActive    *bool  // nil 表示不过滤
}

// RegionRepository 区域仓储接口
type RegionRepository interface {
	// FindTree 根据查询条件返回区域列表（扁平化结构，调用方按需构建树形）
	FindTree(ctx context.Context, db *gorm.DB, query RegionQuery) ([]*Region, error)
}