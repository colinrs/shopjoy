package user

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// RecordOperationLogInput captures all fields needed to write one log row.
// OperatorIP and OperatorUA may be empty if not available from context.
type RecordOperationLogInput struct {
	TenantID     shared.TenantID
	UserID       int64
	Action       string
	OperatorID   int64
	OperatorName string
	Reason       string
	IPAddress    string
	UserAgent    string
}

// OperationLogListItem is the read DTO consumed by the list endpoint.
type OperationLogListItem struct {
	ID           int64
	UserID       int64
	Action       string
	ActionText   string
	OperatorID   int64
	OperatorName string
	Reason       string
	IPAddress    string
	UserAgent    string
	CreatedAt    string
}

// OperationLogListResp is the list response wrapper.
type OperationLogListResp struct {
	List     []*OperationLogListItem
	Total    int64
	Page     int
	PageSize int
}

// OperationLogService persists and reads operation logs.
type OperationLogService interface {
	// Record writes one log row. NEVER returns an error to the caller — the
	// caller MUST treat all returned errors as informational and continue
	// serving the parent business operation. (Errors are logged internally.)
	Record(ctx context.Context, db *gorm.DB, input RecordOperationLogInput)
	List(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, query user.OperationLogQuery) (*OperationLogListResp, error)
}

type operationLogServiceImpl struct {
	db    *gorm.DB
	repo  user.OperationLogRepository
	idGen snowflake.Snowflake
}

// NewOperationLogService wires the application service.
func NewOperationLogService(db *gorm.DB, repo user.OperationLogRepository, idGen snowflake.Snowflake) OperationLogService {
	return &operationLogServiceImpl{db: db, repo: repo, idGen: idGen}
}

func (s *operationLogServiceImpl) Record(ctx context.Context, db *gorm.DB, input RecordOperationLogInput) {
	if db == nil {
		db = s.db
	}
	id, err := s.idGen.NextID(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("operation log: id generation failed: %v", err)
		return
	}
	now := time.Now().UTC()
	entity := &user.OperationLog{
		Model:        application.Model{ID: id, CreatedAt: now, UpdatedAt: now},
		TenantID:     input.TenantID,
		UserID:       input.UserID,
		Action:       input.Action,
		OperatorID:   input.OperatorID,
		OperatorName: input.OperatorName,
		Reason:       input.Reason,
		IPAddress:    input.IPAddress,
		UserAgent:    input.UserAgent,
		Audit:        shared.AuditInfo{CreatedAt: now, UpdatedAt: now},
	}
	if err := s.repo.Create(ctx, db, entity); err != nil {
		logx.WithContext(ctx).Errorf("operation log: write failed (action=%s, user_id=%d): %v", input.Action, input.UserID, err)
	}
}

func (s *operationLogServiceImpl) List(
	ctx context.Context,
	db *gorm.DB,
	tenantID shared.TenantID,
	userID int64,
	query user.OperationLogQuery,
) (*OperationLogListResp, error) {
	if db == nil {
		db = s.db
	}
	logs, total, err := s.repo.FindByUserID(ctx, db, tenantID, userID, query)
	if err != nil {
		return nil, err
	}
	items := make([]*OperationLogListItem, 0, len(logs))
	for _, l := range logs {
		items = append(items, toOperationLogItem(l))
	}
	return &OperationLogListResp{
		List:     items,
		Total:    total,
		Page:     query.Page,
		PageSize: query.Limit(),
	}, nil
}

func toOperationLogItem(l *user.OperationLog) *OperationLogListItem {
	return &OperationLogListItem{
		ID:           l.ID,
		UserID:       l.UserID,
		Action:       l.Action,
		ActionText:   user.ActionText(l.Action),
		OperatorID:   l.OperatorID,
		OperatorName: l.OperatorName,
		Reason:       l.Reason,
		IPAddress:    l.IPAddress,
		UserAgent:    l.UserAgent,
		CreatedAt:    l.Audit.CreatedAt.Format(time.RFC3339),
	}
}