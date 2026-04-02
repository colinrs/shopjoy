package points

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type pointsAccountRepo struct{}

func NewPointsAccountRepository() points.PointsAccountRepository {
	return &pointsAccountRepo{}
}

// Create inserts a new points account
func (r *pointsAccountRepo) Create(ctx context.Context, db *gorm.DB, account *points.PointsAccount) error {
	return db.WithContext(ctx).Create(account).Error
}

// Update updates an existing points account
func (r *pointsAccountRepo) Update(ctx context.Context, db *gorm.DB, account *points.PointsAccount) error {
	return db.WithContext(ctx).
		Model(&points.PointsAccount{}).
		Where("id = ? AND tenant_id = ?", account.ID, account.TenantID.Int64()).
		Save(account).Error
}

// FindByID finds a points account by ID
func (r *pointsAccountRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*points.PointsAccount, error) {
	var account points.PointsAccount
	err := db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPointsAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}

// FindByUserID finds a points account by user ID
func (r *pointsAccountRepo) FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) (*points.PointsAccount, error) {
	var account points.PointsAccount
	err := db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ?", tenantID.Int64(), userID).
		First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPointsAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}

// FindList finds points accounts with pagination and filters
func (r *pointsAccountRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query points.PointsAccountQuery) ([]*points.PointsAccount, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&points.PointsAccount{}).Where("tenant_id = ?", tenantID.Int64())

	if query.UserID != 0 {
		dbQuery = dbQuery.Where("user_id = ?", query.UserID)
	}
	// Note: Email filtering would require joining with users table
	// For simplicity, we skip email filtering here as it would need a join

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var accounts []*points.PointsAccount
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&accounts).Error
	if err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

// GetStats gets statistics for points accounts
func (r *pointsAccountRepo) GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*points.PointsAccountStats, error) {
	stats := &points.PointsAccountStats{}

	// Total accounts
	if err := db.WithContext(ctx).Model(&points.PointsAccount{}).
		Where("tenant_id = ?", tenantID.Int64()).
		Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Total balance (sum of all account balances)
	if err := db.WithContext(ctx).Model(&points.PointsAccount{}).
		Where("tenant_id = ?", tenantID.Int64()).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&stats.TotalBalance).Error; err != nil {
		return nil, err
	}

	// Active accounts (accounts with balance > 0)
	if err := db.WithContext(ctx).Model(&points.PointsAccount{}).
		Where("tenant_id = ? AND balance > 0", tenantID.Int64()).
		Count(&stats.Active).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
