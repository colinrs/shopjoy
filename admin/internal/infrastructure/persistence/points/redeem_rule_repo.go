package points

import (
	"context"
	"errors"
	"fmt"

	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type redeemRuleRepo struct{}

func NewRedeemRuleRepository() points.RedeemRuleRepository {
	return &redeemRuleRepo{}
}

// Create inserts a new redeem rule
func (r *redeemRuleRepo) Create(ctx context.Context, db *gorm.DB, rule *points.RedeemRule) error {
	return db.WithContext(ctx).Create(rule).Error
}

// Update updates an existing redeem rule
func (r *redeemRuleRepo) Update(ctx context.Context, db *gorm.DB, rule *points.RedeemRule) error {
	return db.WithContext(ctx).
		Model(&points.RedeemRule{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", rule.ID, rule.TenantID.Int64()).
		Save(rule).Error
}

// Delete soft deletes a redeem rule
func (r *redeemRuleRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	query := db.WithContext(ctx).Model(&points.RedeemRule{}).Where("id = ? AND deleted_at IS NULL", id)
	// Platform admin (tenantID == 0) can delete all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	result := query.Delete(&points.RedeemRule{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrRedeemRuleNotFound
	}
	return nil
}

// FindByID finds a redeem rule by ID
func (r *redeemRuleRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*points.RedeemRule, error) {
	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var rule points.RedeemRule
	err := query.First(&rule, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrRedeemRuleNotFound
		}
		return nil, err
	}
	return &rule, nil
}

// FindList finds redeem rules with pagination and filters
func (r *redeemRuleRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query points.RedeemRuleQuery) ([]*points.RedeemRule, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&points.RedeemRule{}).Where("deleted_at IS NULL")

	// Tenant filter: Platform admin (TenantID == 0) can access all tenant data
	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.Status != 0 && query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rules []*points.RedeemRule
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&rules).Error
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

// UpdateStatus updates the status of a redeem rule
func (r *redeemRuleRepo) UpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, status points.RedeemRuleStatus) error {
	query := db.WithContext(ctx).
		Model(&points.RedeemRule{}).
		Where("id = ? AND deleted_at IS NULL", id)

	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	result := query.Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrRedeemRuleNotFound
	}
	return nil
}

// IncrementUsedStock atomically increments the used_stock of a redeem rule
func (r *redeemRuleRepo) IncrementUsedStock(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, quantity int64) error {
	query := db.WithContext(ctx).
		Model(&points.RedeemRule{}).
		Where("id = ? AND deleted_at IS NULL", id)

	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	result := query.Update("used_stock", gorm.Expr("used_stock + ?", quantity))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrRedeemRuleNotFound
	}
	return nil
}

// GetStats gets statistics for redeem rules
func (r *redeemRuleRepo) GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*points.RedeemRuleStats, error) {
	stats := &points.RedeemRuleStats{}

	baseQuery := db.WithContext(ctx).Model(&points.RedeemRule{}).Where("deleted_at IS NULL")
	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		baseQuery = baseQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	// Total rules
	if err := baseQuery.Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Active rules
	activeQuery := db.WithContext(ctx).Model(&points.RedeemRule{}).
		Where("status = ? AND deleted_at IS NULL", points.RedeemRuleStatusActive)
	if tenantID != 0 {
		activeQuery = activeQuery.Where("tenant_id = ?", tenantID.Int64())
	}
	if err := activeQuery.Count(&stats.Active).Error; err != nil {
		return nil, err
	}

	// Total redeemed (sum of used_stock)
	sumQuery := db.WithContext(ctx).Model(&points.RedeemRule{}).
		Where("deleted_at IS NULL")
	if tenantID != 0 {
		sumQuery = sumQuery.Where("tenant_id = ?", tenantID.Int64())
	}
	if err := sumQuery.Select("COALESCE(SUM(used_stock), 0)").Scan(&stats.TotalRedeemed).Error; err != nil {
		return nil, err
	}

	return stats, nil
}