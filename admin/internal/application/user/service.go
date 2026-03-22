package user

import (
	"context"

	domain "github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

// DTOs
type CreateUserRequest struct {
	TenantID shared.TenantID
	Email    string
	Phone    string
	Password string
	Name     string
	Avatar   string
	Gender   domain.Gender
	Birthday string
}

type UpdateUserRequest struct {
	ID       int64
	TenantID shared.TenantID
	Name     string
	Avatar   string
	Gender   domain.Gender
	Birthday string
}

type ChangePasswordRequest struct {
	UserID          int64
	TenantID        shared.TenantID
	OldPassword     string
	NewPassword     string
	ConfirmPassword string
}

type UserResponse struct {
	ID        int64
	TenantID  int64
	Email     string
	Phone     string
	Name      string
	Avatar    string
	Gender    int
	Birthday  string
	Status    int
	LastLogin string
	CreatedAt string
}

type UserListResponse struct {
	List     []*UserResponse
	Total    int64
	Page     int
	PageSize int
}

type QueryRequest struct {
	shared.PageQuery
	Name   string
	Email  string
	Phone  string
	Status domain.Status
}

type UserStatsResponse struct {
	Total     int64
	Active    int64
	Suspended int64
	NewToday  int64
}

// Service 用户应用服务接口
type Service interface {
	Register(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
	Update(ctx context.Context, req UpdateUserRequest) (*UserResponse, error)
	ChangePassword(ctx context.Context, req ChangePasswordRequest) error
	GetByID(ctx context.Context, tenantID shared.TenantID, id int64) (*UserResponse, error)
	GetByEmail(ctx context.Context, tenantID shared.TenantID, email string) (*UserResponse, error)
	List(ctx context.Context, tenantID shared.TenantID, req QueryRequest) (*UserListResponse, error)
	Suspend(ctx context.Context, tenantID shared.TenantID, id int64) error
	Activate(ctx context.Context, tenantID shared.TenantID, id int64) error
	Delete(ctx context.Context, tenantID shared.TenantID, id int64) error
	ResetPassword(ctx context.Context, tenantID shared.TenantID, id int64) (string, error)
	GetStats(ctx context.Context, tenantID shared.TenantID) (*UserStatsResponse, error)
}
