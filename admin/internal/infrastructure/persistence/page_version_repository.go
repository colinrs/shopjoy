package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/storefront"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type pageVersionRepo struct{}

func NewPageVersionRepository() storefront.PageVersionRepository {
	return &pageVersionRepo{}
}

type pageVersionModel struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	TenantID  int64     `gorm:"column:tenant_id;not null;uniqueIndex:idx_tenant_page_ver"`
	PageID    int64     `gorm:"column:page_id;not null;uniqueIndex:idx_tenant_page_ver;index"`
	Version   int       `gorm:"column:version;not null;uniqueIndex:idx_tenant_page_ver"`
	Blocks    string    `gorm:"column:blocks;type:text;not null"`
	CreatedBy int64     `gorm:"column:created_by;not null;default:0"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (pageVersionModel) TableName() string {
	return "page_versions"
}

func (m *pageVersionModel) toEntity() *storefront.PageVersion {
	var blocks []storefront.BlockSnapshot
	if m.Blocks != "" {
		json.Unmarshal([]byte(m.Blocks), &blocks)
	}
	return &storefront.PageVersion{
		Model:     application.Model{ID: m.ID},
		TenantID:  shared.TenantID(m.TenantID),
		PageID:    m.PageID,
		Version:   m.Version,
		Blocks:    blocks,
		CreatedBy: m.CreatedBy,
	}
}

func fromPageVersionEntity(v *storefront.PageVersion) *pageVersionModel {
	blocks, _ := json.Marshal(v.Blocks)
	return &pageVersionModel{
		ID:        v.ID,
		TenantID:  v.TenantID.Int64(),
		PageID:    v.PageID,
		Version:   v.Version,
		Blocks:    string(blocks),
		CreatedBy: v.CreatedBy,
		CreatedAt: time.Now().UTC(),
	}
}

func (r *pageVersionRepo) Create(ctx context.Context, db *gorm.DB, v *storefront.PageVersion) error {
	model := fromPageVersionEntity(v)
	return db.WithContext(ctx).Create(model).Error
}

func (r *pageVersionRepo) FindByPageID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageID int64, page, pageSize int) ([]*storefront.PageVersion, int64, error) {
	var total int64
	if err := db.WithContext(ctx).Model(&pageVersionModel{}).
		Where("page_id = ? AND tenant_id = ?", pageID, tenantID.Int64()).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []pageVersionModel
	offset := (page - 1) * pageSize
	err := db.WithContext(ctx).
		Where("page_id = ? AND tenant_id = ?", pageID, tenantID.Int64()).
		Order("version DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	versions := make([]*storefront.PageVersion, len(models))
	for i, m := range models {
		versions[i] = m.toEntity()
	}
	return versions, total, nil
}

func (r *pageVersionRepo) CountByPageID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageID int64) (int64, error) {
	var total int64
	err := db.WithContext(ctx).Model(&pageVersionModel{}).
		Where("page_id = ? AND tenant_id = ?", pageID, tenantID.Int64()).
		Count(&total).Error
	return total, err
}

func (r *pageVersionRepo) FindByVersion(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageID int64, version int) (*storefront.PageVersion, error) {
	var model pageVersionModel
	err := db.WithContext(ctx).
		Where("page_id = ? AND tenant_id = ? AND version = ?", pageID, tenantID.Int64(), version).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *pageVersionRepo) DeleteOldest(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageID int64, keepCount int) error {
	// Get total count
	var total int64
	if err := db.WithContext(ctx).Model(&pageVersionModel{}).
		Where("page_id = ? AND tenant_id = ?", pageID, tenantID.Int64()).
		Count(&total).Error; err != nil {
		return err
	}

	if int(total) <= keepCount {
		return nil
	}

	// Find IDs to delete (oldest versions)
	var idsToDelete []int64
	err := db.WithContext(ctx).Model(&pageVersionModel{}).
		Select("id").
		Where("page_id = ? AND tenant_id = ?", pageID, tenantID.Int64()).
		Order("version ASC").
		Limit(int(total)-keepCount).
		Pluck("id", &idsToDelete).Error
	if err != nil {
		return err
	}

	if len(idsToDelete) == 0 {
		return nil
	}

	return db.WithContext(ctx).
		Where("id IN ?", idsToDelete).
		Delete(&pageVersionModel{}).Error
}
