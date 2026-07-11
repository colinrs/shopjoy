package media

import (
	"time"

	"github.com/colinrs/shopjoy/pkg/application"
)

// Asset 媒体资产元数据（与 sql/admin/schema.sql 的 media_assets 表一一对应）
type Asset struct {
	application.Model
	PublicID  string `gorm:"column:public_id;size:255;not null;uniqueIndex:uk_provider_public"`
	URL       string `gorm:"column:url;size:1024;not null"`
	Filename  string `gorm:"column:filename;size:255;not null;default:''"`
	SizeBytes int64  `gorm:"column:size_bytes;not null;default:0"`
	MimeType  string `gorm:"column:mime_type;size:64;not null;default:''"`
	Width     int    `gorm:"column:width;not null;default:0"`
	Height    int    `gorm:"column:height;not null;default:0"`
	Format    string `gorm:"column:format;size:32;not null;default:''"`
	Category  string `gorm:"column:category;size:32;not null;default:'common';index:idx_tenant_category"`
	Provider  string `gorm:"column:provider;size:16;not null;default:'local';uniqueIndex:uk_provider_public"`
	TenantID  int64  `gorm:"column:tenant_id;not null;index:idx_tenant_category"`
	CreatedBy int64  `gorm:"column:created_by;not null;default:0"`
}

// TableName implements gorm.Tabler.
func (*Asset) TableName() string { return "media_assets" }

// Tenant returns the asset's tenant ID. application.Model.ID is int64.
func (a *Asset) Tenant() int64 { return a.TenantID }

// SizeInBytes returns the asset's size in bytes.
func (a *Asset) SizeInBytes() int64 { return a.SizeBytes }

// Now returns current UTC time — re-export to avoid importing "time" in tests later.
func Now() time.Time { return time.Now().UTC() }