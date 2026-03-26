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

type earnRuleRepo struct{}

func NewEarnRuleRepository() points.EarnRuleRepository {
	return &earnRuleRepo{}
}

// Create inserts a new earn rule
func (r *earnRuleRepo) Create(ctx context.Context, db *gorm.DB, rule *points.EarnRule) error {
	return db.WithContext(ctx).Create(rule).Error
}

// Update updates an existing earn rule
func (r *earnRuleRepo) Update(ctx context.Context, db *gorm.DB, rule *points.EarnRule) error {
	return db.WithContext(ctx).
		Model(&points.EarnRule{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", rule.ID, rule.TenantID.Int64()).
		Save(rule).Error
}

// Delete soft deletes an earn rule
func (r *earnRuleRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	query := db.WithContext(ctx).Model(&points.EarnRule{}).Where("id = ? AND deleted_at IS NULL", id)
	// Platform admin (tenantID == 0) can delete all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	result := query.Delete(&points.EarnRule{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrEarnRuleNotFound
	}
	return nil
}

// FindByID finds an earn rule by ID
func (r *earnRuleRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*points.EarnRule, error) {
	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var rule points.EarnRule
	err := query.First(&rule, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrEarnRuleNotFound
		}
		return nil, err
	}
	return &rule, nil
}

// FindList finds earn rules with pagination and filters
func (r *earnRuleRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query points.EarnRuleQuery) ([]*points.EarnRule, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&points.EarnRule{}).Where("deleted_at IS NULL")

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
	if query.Scenario != "" && query.Scenario.IsValid() {
		dbQuery = dbQuery.Where("scenario = ?", query.Scenario)
	}
	if query.CalculationType != "" && query.CalculationType.IsValid() {
		dbQuery = dbQuery.Where("calculation_type = ?", query.CalculationType)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rules []*points.EarnRule
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&rules).Error
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

// FindByScenario finds earn rules by scenario
func (r *earnRuleRepo) FindByScenario(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, scenario points.EarnScenario) ([]*points.EarnRule, error) {
	query := db.WithContext(ctx).Model(&points.EarnRule{}).
		Where("scenario = ? AND deleted_at IS NULL", scenario)

	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	var rules []*points.EarnRule
	err := query.Order("priority DESC, created_at DESC").Find(&rules).Error
	if err != nil {
		return nil, err
	}

	return rules, nil
}

// UpdateStatus updates the status of an earn rule
func (r *earnRuleRepo) UpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, status points.EarnRuleStatus) error {
	query := db.WithContext(ctx).
		Model(&points.EarnRule{}).
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
		return code.ErrEarnRuleNotFound
	}
	return nil
}

// GetStats gets statistics for earn rules
func (r *earnRuleRepo) GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*points.EarnRuleStats, error) {
	stats := &points.EarnRuleStats{}

	baseQuery := db.WithContext(ctx).Model(&points.EarnRule{}).Where("deleted_at IS NULL")
	// Platform admin (tenantID == 0) can access all tenant data
	if tenantID != 0 {
		baseQuery = baseQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	// Total rules
	if err := baseQuery.Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Active rules
	activeQuery := db.WithContext(ctx).Model(&points.EarnRule{}).
		Where("status = ? AND deleted_at IS NULL", points.EarnRuleStatusActive)
	if tenantID != 0 {
		activeQuery = activeQuery.Where("tenant_id = ?", tenantID.Int64())
	}
	if err := activeQuery.Count(&stats.Active).Error; err != nil {
		return nil, err
	}

	return stats, nil
}