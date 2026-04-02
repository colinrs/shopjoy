package points

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type pointsRedemptionRepo struct{}

func NewPointsRedemptionRepository() points.PointsRedemptionRepository {
	return &pointsRedemptionRepo{}
}

// Create inserts a new points redemption
func (r *pointsRedemptionRepo) Create(ctx context.Context, db *gorm.DB, redemption *points.PointsRedemption) error {
	return db.WithContext(ctx).Create(redemption).Error
}

// Update updates an existing points redemption
func (r *pointsRedemptionRepo) Update(ctx context.Context, db *gorm.DB, redemption *points.PointsRedemption) error {
	return db.WithContext(ctx).
		Model(&points.PointsRedemption{}).
		Where("id = ? AND tenant_id = ?", redemption.ID, redemption.TenantID.Int64()).
		Save(redemption).Error
}

// FindByID finds a points redemption by ID
func (r *pointsRedemptionRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*points.PointsRedemption, error) {
	var redemption points.PointsRedemption
	err := db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		First(&redemption).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPointsRedemptionNotFound
		}
		return nil, err
	}
	return &redemption, nil
}

// FindList finds points redemptions with pagination and filters
func (r *pointsRedemptionRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query points.PointsRedemptionQuery) ([]*points.PointsRedemption, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&points.PointsRedemption{}).Where("tenant_id = ?", tenantID.Int64())

	if query.UserID != 0 {
		dbQuery = dbQuery.Where("user_id = ?", query.UserID)
	}
	if query.Status != 0 && query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.StartTime != nil {
		dbQuery = dbQuery.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime != nil {
		dbQuery = dbQuery.Where("created_at < ?", query.EndTime)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var redemptions []*points.PointsRedemption
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&redemptions).Error
	if err != nil {
		return nil, 0, err
	}

	return redemptions, total, nil
}

// CountByUserAndRule counts the number of redemptions by a user for a specific rule
func (r *pointsRedemptionRepo) CountByUserAndRule(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID, ruleID int64) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Model(&points.PointsRedemption{}).
		Where("tenant_id = ? AND user_id = ? AND redeem_rule_id = ? AND status != ?",
			tenantID.Int64(), userID, ruleID, points.RedemptionStatusCancelled).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
