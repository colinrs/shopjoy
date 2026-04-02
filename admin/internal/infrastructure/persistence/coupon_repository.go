package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type couponRepo struct{}

func NewCouponRepository() promotion.CouponRepository {
	return &couponRepo{}
}

// couponModel represents the database model for Coupon
type couponModel struct {
	ID           int64           `gorm:"column:id;primaryKey;autoIncrement:false"`
	TenantID     int64           `gorm:"column:tenant_id;not null;index"`
	Name         string          `gorm:"column:name;size:200;not null;index"`
	Code         string          `gorm:"column:code;size:50;not null;uniqueIndex"`
	Description  string          `gorm:"column:description;type:text"`
	Type         int             `gorm:"column:type;not null;index"`
	Value        decimal.Decimal `gorm:"column:value;type:decimal(19,4);not null"`
	MinAmount    decimal.Decimal `gorm:"column:min_amount;type:decimal(19,4);not null;default:0"`
	MaxDiscount  decimal.Decimal `gorm:"column:max_discount;type:decimal(19,4);not null;default:0"`
	Currency     string          `gorm:"column:currency;size:10;not null"`
	TotalCount   int             `gorm:"column:total_count;not null;default:0"`
	UsedCount    int             `gorm:"column:used_count;not null;default:0"`
	PerUserLimit int             `gorm:"column:per_user_limit;not null;default:0"`
	Status       int             `gorm:"column:status;not null;index"`
	StartAt      int64           `gorm:"column:start_at;not null;index"`
	EndAt        int64           `gorm:"column:end_at;not null;index"`
	ScopeType    string          `gorm:"column:scope_type;size:32;not null"`
	ScopeIDs     string          `gorm:"column:scope_ids;type:json"`   // JSON array of int64
	ExcludeIDs   string          `gorm:"column:exclude_ids;type:json"` // JSON array of int64
	CreatedBy    int64           `gorm:"column:created_by;not null"`
	UpdatedBy    int64           `gorm:"column:updated_by;not null"`
	DeletedAt    *int64          `gorm:"column:deleted_at;index"`
	CreatedAt    int64           `gorm:"column:created_at"`
	UpdatedAt    int64           `gorm:"column:updated_at"`
}

func (couponModel) TableName() string {
	return "coupons"
}

func (m *couponModel) toEntity() *promotion.Coupon {
	// Parse JSON arrays for scope
	var scopeIDs, excludeIDs []int64
	if m.ScopeIDs != "" {
		_ = json.Unmarshal([]byte(m.ScopeIDs), &scopeIDs)
	}
	if m.ExcludeIDs != "" {
		_ = json.Unmarshal([]byte(m.ExcludeIDs), &excludeIDs)
	}

	return &promotion.Coupon{
		ID:           m.ID,
		TenantID:     shared.TenantID(m.TenantID),
		Name:         m.Name,
		Code:         m.Code,
		Description:  m.Description,
		Type:         promotion.CouponType(m.Type),
		Value:        m.Value,
		MinAmount:    m.MinAmount,
		MaxDiscount:  m.MaxDiscount,
		Currency:     m.Currency,
		TotalCount:   m.TotalCount,
		UsedCount:    m.UsedCount,
		PerUserLimit: m.PerUserLimit,
		Status:       promotion.CouponStatus(m.Status),
		StartAt:      time.Unix(m.StartAt, 0),
		EndAt:        time.Unix(m.EndAt, 0),
		Scope: promotion.PromotionScope{
			Type:       promotion.ScopeType(m.ScopeType),
			IDs:        scopeIDs,
			ExcludeIDs: excludeIDs,
		},
		Audit: shared.AuditInfo{
			CreatedAt: time.Unix(m.CreatedAt, 0).UTC(),
			UpdatedAt: time.Unix(m.UpdatedAt, 0).UTC(),
			CreatedBy: m.CreatedBy,
			UpdatedBy: m.UpdatedBy,
		},
		DeletedAt: m.DeletedAt,
	}
}

func fromCouponEntity(c *promotion.Coupon) *couponModel {
	// Serialize JSON arrays for scope
	scopeIDsJSON, _ := json.Marshal(c.Scope.IDs)
	excludeIDsJSON, _ := json.Marshal(c.Scope.ExcludeIDs)

	return &couponModel{
		ID:           c.ID,
		TenantID:     c.TenantID.Int64(),
		Name:         c.Name,
		Code:         c.Code,
		Description:  c.Description,
		Type:         int(c.Type),
		Value:        c.Value,
		MinAmount:    c.MinAmount,
		MaxDiscount:  c.MaxDiscount,
		Currency:     c.Currency,
		TotalCount:   c.TotalCount,
		UsedCount:    c.UsedCount,
		PerUserLimit: c.PerUserLimit,
		Status:       int(c.Status),
		StartAt:      c.StartAt.Unix(),
		EndAt:        c.EndAt.Unix(),
		ScopeType:    string(c.Scope.Type),
		ScopeIDs:     string(scopeIDsJSON),
		ExcludeIDs:   string(excludeIDsJSON),
		CreatedBy:    c.Audit.CreatedBy,
		UpdatedBy:    c.Audit.UpdatedBy,
		DeletedAt:    c.DeletedAt,
		CreatedAt:    c.Audit.CreatedAt.Unix(),
		UpdatedAt:    c.Audit.UpdatedAt.Unix(),
	}
}

// Create inserts a new coupon
func (r *couponRepo) Create(ctx context.Context, db *gorm.DB, c *promotion.Coupon) error {
	model := fromCouponEntity(c)
	return db.WithContext(ctx).Create(model).Error
}

// Update updates an existing coupon
func (r *couponRepo) Update(ctx context.Context, db *gorm.DB, c *promotion.Coupon) error {
	model := fromCouponEntity(c)
	return db.WithContext(ctx).
		Model(&couponModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", c.ID, c.TenantID.Int64()).
		Updates(map[string]interface{}{
			"name":           model.Name,
			"code":           model.Code,
			"description":    model.Description,
			"type":           model.Type,
			"value":          model.Value,
			"min_amount":     model.MinAmount,
			"max_discount":   model.MaxDiscount,
			"currency":       model.Currency,
			"total_count":    model.TotalCount,
			"used_count":     model.UsedCount,
			"per_user_limit": model.PerUserLimit,
			"status":         model.Status,
			"start_at":       model.StartAt,
			"end_at":         model.EndAt,
			"scope_type":     model.ScopeType,
			"scope_ids":      model.ScopeIDs,
			"exclude_ids":    model.ExcludeIDs,
			"updated_by":     model.UpdatedBy,
			"updated_at":     model.UpdatedAt,
		}).Error
}

// Delete soft deletes a coupon
func (r *couponRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	query := db.WithContext(ctx).Model(&couponModel{}).Where("id = ? AND deleted_at IS NULL", id)
	// Platform admin (tenantID == 0) can delete all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	now := time.Now().UTC()
	result := query.Update("deleted_at", now)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrCouponNotFound
	}
	return nil
}

// FindByID finds a coupon by ID
func (r *couponRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*promotion.Coupon, error) {
	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model couponModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrCouponNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindByCode finds a coupon by code
func (r *couponRepo) FindByCode(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, codeStr string) (*promotion.Coupon, error) {
	query := db.WithContext(ctx).Model(&couponModel{}).Where("code = ? AND deleted_at IS NULL", codeStr)
	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model couponModel
	err := query.First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrCouponNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindList finds coupons with pagination and filters
func (r *couponRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query promotion.CouponQuery) ([]*promotion.Coupon, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&couponModel{}).Where("deleted_at IS NULL")

	// Tenant filter: Platform admin (TenantID == 0) can access all tenant data
	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.Code != "" {
		dbQuery = dbQuery.Where("code LIKE ?", fmt.Sprintf("%%%s%%", query.Code))
	}
	if query.Status != nil && query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", *query.Status)
	}
	if query.Type != nil && query.Type.IsValid() {
		dbQuery = dbQuery.Where("type = ?", *query.Type)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []couponModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	coupons := make([]*promotion.Coupon, len(models))
	for i, m := range models {
		coupons[i] = m.toEntity()
	}
	return coupons, total, nil
}

// IncrementUsage atomically increments the used_count of a coupon
func (r *couponRepo) IncrementUsage(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	query := db.WithContext(ctx).
		Model(&couponModel{}).
		Where("id = ? AND deleted_at IS NULL", id)

	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	result := query.Update("used_count", gorm.Expr("used_count + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrCouponNotFound
	}
	return nil
}
