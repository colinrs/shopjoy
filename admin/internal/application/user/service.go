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

// AddressResponse 地址响应
type AddressResponse struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Country    string `json:"country"`
	Province   string `json:"province"`
	City       string `json:"city"`
	District   string `json:"district"`
	Detail     string `json:"detail"`
	PostalCode string `json:"postal_code"`
	IsDefault  bool   `json:"is_default"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type AddressListResponse struct {
	List  []*AddressResponse `json:"list"`
	Total int64              `json:"total"`
}

type UserDetailResponse struct {
	ID            int64  `json:"id"`
	TenantID      int64  `json:"tenant_id"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Name          string `json:"name"`
	Avatar        string `json:"avatar"`
	Gender        int    `json:"gender"`
	GenderText    string `json:"gender_text"`
	Birthday      string `json:"birthday"`
	Status        int    `json:"status"`
	StatusText    string `json:"status_text"`
	PointsBalance int64  `json:"points_balance"`
	PointsFrozen  int64  `json:"points_frozen"`
	TotalEarned   int64  `json:"total_earned"`
	TotalRedeemed int64  `json:"total_redeemed"`
	OrderCount    int64  `json:"order_count"`
	TotalSpent    string `json:"total_spent"`
	ReviewCount   int64  `json:"review_count"`
	LastLogin     string `json:"last_login"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type EnhancedQueryRequest struct {
	shared.PageQuery
	Keyword       string
	Status        domain.Status
	RegisterStart string
	RegisterEnd   string
}

type ExtendedListResponse struct {
	List     []*ExtendedUserResponse `json:"list"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
}

type ExtendedUserResponse struct {
	ID            int64  `json:"id"`
	TenantID      int64  `json:"tenant_id"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Name          string `json:"name"`
	Avatar        string `json:"avatar"`
	Status        int    `json:"status"`
	StatusText    string `json:"status_text"`
	PointsBalance int64  `json:"points_balance"`
	OrderCount    int64  `json:"order_count"`
	TotalSpent    string `json:"total_spent"`
	LastLogin     string `json:"last_login"`
	CreatedAt     string `json:"created_at"`
}

type SuspendUserRequest struct {
	TenantID shared.TenantID
	UserID   int64
	Reason   string
}

type UserStats struct {
	TotalUsers     int64 `json:"total_users"`
	ActiveUsers    int64 `json:"active_users"`
	SuspendedUsers int64 `json:"suspended_users"`
	NewUsersToday  int64 `json:"new_users_today"`
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
	// Extended methods
	GetDetail(ctx context.Context, tenantID shared.TenantID, id int64) (*UserDetailResponse, error)
	ExtendedList(ctx context.Context, tenantID shared.TenantID, req EnhancedQueryRequest) (*ExtendedListResponse, error)
	GetAddresses(ctx context.Context, tenantID shared.TenantID, userID int64) (*AddressListResponse, error)
	SuspendWithReason(ctx context.Context, req SuspendUserRequest) error
	ActivateUser(ctx context.Context, tenantID shared.TenantID, userID int64) error
	DeleteUser(ctx context.Context, tenantID shared.TenantID, userID int64) error
	GetUserStats(ctx context.Context, tenantID shared.TenantID) (*UserStats, error)
	ExportUsers(ctx context.Context, tenantID shared.TenantID, req EnhancedQueryRequest) ([]byte, error)
}
