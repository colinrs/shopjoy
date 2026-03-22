// Package product 商品领域层
package product

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// ProductLocalization 商品多语言本地化实体
type ProductLocalization struct {
	ID           int64           // 本地化ID
	TenantID     shared.TenantID // 租户ID
	ProductID    int64           // 商品ID
	LanguageCode string          // 语言代码，如 en, zh-CN, ja
	Name         string          // 商品名称（本地化）
	Description  string          // 商品描述（本地化）
	AuditInfo    shared.AuditInfo `gorm:"embedded"` // 审计信息
}

// TableName 表名
func (p *ProductLocalization) TableName() string {
	return "product_localizations"
}

// ProductLocalizationRepository 商品本地化仓储接口
type ProductLocalizationRepository interface {
	// Create 创建商品本地化记录
	Create(ctx context.Context, db *gorm.DB, localization *ProductLocalization) error
	// Update 更新商品本地化记录
	Update(ctx context.Context, db *gorm.DB, localization *ProductLocalization) error
	// Delete 删除商品本地化记录
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	// FindByID 根据ID查找商品本地化记录
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*ProductLocalization, error)
	// FindByProductID 根据商品ID查找所有本地化记录
	FindByProductID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, productID int64) ([]*ProductLocalization, error)
	// FindByProductAndLanguage 根据商品ID和语言代码查找本地化记录
	FindByProductAndLanguage(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, productID int64, languageCode string) (*ProductLocalization, error)
}