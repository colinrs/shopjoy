package user

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// Action constants for user operation logs.
// Single source of truth — frontend mirrors via i18n keys.
const (
	ActionCreateUser        = "CREATE_USER"
	ActionUpdateUser        = "UPDATE_USER"
	ActionSuspendUser       = "SUSPEND_USER"
	ActionSuspendWithReason = "SUSPEND_WITH_REASON"
	ActionActivateUser      = "ACTIVATE_USER"
	ActionDeleteUser        = "DELETE_USER"
	ActionResetPassword     = "RESET_PASSWORD"
)

// OperationLog records an admin-side action against a user (state changes,
// profile updates, password resets). Writes are best-effort: instrumentation
// failures must never block the parent business operation.
type OperationLog struct {
	application.Model
	TenantID     shared.TenantID
	UserID       int64
	Action       string
	OperatorID   int64
	OperatorName string
	Reason       string
	IPAddress    string
	UserAgent    string
	Audit        shared.AuditInfo `gorm:"embedded"`
}

func (o *OperationLog) TableName() string {
	return "user_operation_logs"
}

// OperationLogQuery filters for FindByUserID.
type OperationLogQuery struct {
	Page     int
	PageSize int
	Action   string // empty = all
	Keyword  string // empty = no keyword filter
}

func (q OperationLogQuery) Offset() int {
	if q.Page <= 0 {
		q.Page = 1
	}
	return (q.Page - 1) * q.PageSize
}

func (q OperationLogQuery) Limit() int {
	if q.PageSize <= 0 {
		return 20
	}
	return q.PageSize
}

// OperationLogRepository persists and queries OperationLog records.
// Read paths respect platform-admin tenant bypass: tenantID == 0 means
// "no tenant filter".
type OperationLogRepository interface {
	Create(ctx context.Context, db *gorm.DB, log *OperationLog) error
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, query OperationLogQuery) ([]*OperationLog, int64, error)
}

// ActionText returns the Chinese display label for a known action.
// Returns the raw action string for unknown values so new actions
// degrade gracefully rather than disappearing.
func ActionText(action string) string {
	switch action {
	case ActionCreateUser:
		return "创建用户"
	case ActionUpdateUser:
		return "更新资料"
	case ActionSuspendUser, ActionSuspendWithReason:
		return "禁用用户"
	case ActionActivateUser:
		return "启用用户"
	case ActionDeleteUser:
		return "删除用户"
	case ActionResetPassword:
		return "重置密码"
	default:
		return action
	}
}
