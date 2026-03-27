package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/storefront"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type themeAuditLogRepo struct{}

func NewThemeAuditLogRepository() storefront.ThemeAuditLogRepository {
	return &themeAuditLogRepo{}
}

type themeAuditLogModel struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	TenantID  int64     `gorm:"column:tenant_id;not null;index"`
	Action    string    `gorm:"column:action;type:varchar(30);not null"`
	ThemeID   int64     `gorm:"column:theme_id;not null"`
	ThemeName string    `gorm:"column:theme_name;type:varchar(100);not null"`
	ThemeCode string    `gorm:"column:theme_code;type:varchar(50);not null"`
	OldConfig string    `gorm:"column:old_config;type:text"`
	NewConfig string    `gorm:"column:new_config;type:text"`
	UserID    int64     `gorm:"column:user_id;not null"`
	UserName  string    `gorm:"column:user_name;type:varchar(100);not null;default:''"`
	IPAddress string    `gorm:"column:ip_address;type:varchar(45);default:''"`
	UserAgent string    `gorm:"column:user_agent;type:varchar(500);default:''"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (themeAuditLogModel) TableName() string {
	return "theme_audit_logs"
}

func (m *themeAuditLogModel) toEntity() *storefront.ThemeAuditLog {
	return &storefront.ThemeAuditLog{
		Model:     application.Model{ID: m.ID},
		TenantID:  shared.TenantID(m.TenantID),
		Action:    m.Action,
		ThemeID:   m.ThemeID,
		ThemeName: m.ThemeName,
		ThemeCode: m.ThemeCode,
		OldConfig: m.OldConfig,
		NewConfig: m.NewConfig,
		UserID:    m.UserID,
		UserName:  m.UserName,
		IPAddress: m.IPAddress,
		UserAgent: m.UserAgent,
	}
}

func (r *themeAuditLogRepo) Create(ctx context.Context, db *gorm.DB, log *storefront.ThemeAuditLog) error {
	model := &themeAuditLogModel{
		ID:        log.ID,
		TenantID:  log.TenantID.Int64(),
		Action:    log.Action,
		ThemeID:   log.ThemeID,
		ThemeName: log.ThemeName,
		ThemeCode: log.ThemeCode,
		OldConfig: log.OldConfig,
		NewConfig: log.NewConfig,
		UserID:    log.UserID,
		UserName:  log.UserName,
		IPAddress: log.IPAddress,
		UserAgent: log.UserAgent,
		CreatedAt: time.Now().UTC(),
	}
	return db.WithContext(ctx).Create(model).Error
}

func (r *themeAuditLogRepo) FindByTenantID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, page, pageSize int) ([]*storefront.ThemeAuditLog, int64, error) {
	var total int64
	if err := db.WithContext(ctx).Model(&themeAuditLogModel{}).
		Where("tenant_id = ?", tenantID.Int64()).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []themeAuditLogModel
	offset := (page - 1) * pageSize
	err := db.WithContext(ctx).
		Where("tenant_id = ?", tenantID.Int64()).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	logs := make([]*storefront.ThemeAuditLog, len(models))
	for i, m := range models {
		logs[i] = m.toEntity()
	}
	return logs, total, nil
}