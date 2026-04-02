package points

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type pointsTransactionRepo struct{}

func NewPointsTransactionRepository() points.PointsTransactionRepository {
	return &pointsTransactionRepo{}
}

// Create inserts a new points transaction
func (r *pointsTransactionRepo) Create(ctx context.Context, db *gorm.DB, transaction *points.PointsTransaction) error {
	return db.WithContext(ctx).Create(transaction).Error
}

// FindByID finds a points transaction by ID
func (r *pointsTransactionRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*points.PointsTransaction, error) {
	query := db.WithContext(ctx).Where("tenant_id = ?", tenantID.Int64())
	var transaction points.PointsTransaction
	err := query.First(&transaction, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPointsTransactionNotFound
		}
		return nil, err
	}
	return &transaction, nil
}

// FindList finds points transactions with pagination and filters
func (r *pointsTransactionRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query points.PointsTransactionQuery) ([]*points.PointsTransaction, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&points.PointsTransaction{}).Where("tenant_id = ?", tenantID.Int64())

	if query.UserID != 0 {
		dbQuery = dbQuery.Where("user_id = ?", query.UserID)
	}
	if query.AccountID != 0 {
		dbQuery = dbQuery.Where("account_id = ?", query.AccountID)
	}
	if query.Type != "" && query.Type.IsValid() {
		dbQuery = dbQuery.Where("type = ?", query.Type)
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

	var transactions []*points.PointsTransaction
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetStats gets statistics for points transactions
func (r *pointsTransactionRepo) GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query points.PointsTransactionQuery) (*points.PointsTransactionStats, error) {
	stats := &points.PointsTransactionStats{}

	// Base query with tenant filter
	dbQuery := db.WithContext(ctx).Model(&points.PointsTransaction{}).Where("tenant_id = ?", tenantID.Int64())

	// Apply filters from query
	if query.UserID != 0 {
		dbQuery = dbQuery.Where("user_id = ?", query.UserID)
	}
	if query.AccountID != 0 {
		dbQuery = dbQuery.Where("account_id = ?", query.AccountID)
	}
	if query.Type != "" && query.Type.IsValid() {
		dbQuery = dbQuery.Where("type = ?", query.Type)
	}
	if query.StartTime != nil {
		dbQuery = dbQuery.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime != nil {
		dbQuery = dbQuery.Where("created_at < ?", query.EndTime)
	}

	// Total earned (sum of positive points from EARN transactions)
	earnQuery := db.WithContext(ctx).Model(&points.PointsTransaction{}).
		Where("tenant_id = ?", tenantID.Int64()).
		Where("type = ?", points.TransactionTypeEarn)
	if query.UserID != 0 {
		earnQuery = earnQuery.Where("user_id = ?", query.UserID)
	}
	if query.AccountID != 0 {
		earnQuery = earnQuery.Where("account_id = ?", query.AccountID)
	}
	if query.StartTime != nil {
		earnQuery = earnQuery.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime != nil {
		earnQuery = earnQuery.Where("created_at < ?", query.EndTime)
	}
	if err := earnQuery.Select("COALESCE(SUM(points), 0)").Scan(&stats.TotalEarned).Error; err != nil {
		return nil, err
	}

	// Total redeemed (sum of absolute value of negative points from REDEEM transactions)
	redeemQuery := db.WithContext(ctx).Model(&points.PointsTransaction{}).
		Where("tenant_id = ?", tenantID.Int64()).
		Where("type = ?", points.TransactionTypeRedeem)
	if query.UserID != 0 {
		redeemQuery = redeemQuery.Where("user_id = ?", query.UserID)
	}
	if query.AccountID != 0 {
		redeemQuery = redeemQuery.Where("account_id = ?", query.AccountID)
	}
	if query.StartTime != nil {
		redeemQuery = redeemQuery.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime != nil {
		redeemQuery = redeemQuery.Where("created_at < ?", query.EndTime)
	}
	// Points for REDEEM transactions are negative, so we take ABS
	if err := redeemQuery.Select("COALESCE(SUM(ABS(points)), 0)").Scan(&stats.TotalRedeemed).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
