package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type OperationLogRepositoryImpl struct{}

func NewOperationLogRepository() user.OperationLogRepository {
	return &OperationLogRepositoryImpl{}
}

func (r *OperationLogRepositoryImpl) Create(ctx context.Context, db *gorm.DB, log *user.OperationLog) error {
	now := time.Now().UTC()
	if log.CreatedAt.IsZero() {
		log.CreatedAt = now
	}
	if log.UpdatedAt.IsZero() {
		log.UpdatedAt = now
	}
	// AuditInfo is embedded — populate it so both the column data and the
	// in-memory entity carry the timestamps callers read via l.Audit.CreatedAt.
	if log.Audit.CreatedAt.IsZero() {
		log.Audit.CreatedAt = now
	}
	if log.Audit.UpdatedAt.IsZero() {
		log.Audit.UpdatedAt = now
	}
	return db.WithContext(ctx).Create(log).Error
}

func (r *OperationLogRepositoryImpl) FindByUserID(
	ctx context.Context,
	db *gorm.DB,
	tenantID shared.TenantID,
	userID int64,
	query user.OperationLogQuery,
) ([]*user.OperationLog, int64, error) {
	q := db.WithContext(ctx).Model(&user.OperationLog{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	// Platform admin (tenantID == 0) sees logs across all tenants.
	if tenantID != 0 {
		q = q.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.Action != "" {
		q = q.Where("action = ?", query.Action)
	}
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		q = q.Where("(operator_name LIKE ? OR reason LIKE ?)", like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []*user.OperationLog
	err := q.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&logs).Error
	return logs, total, err
}
