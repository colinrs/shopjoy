package adminuser

import (
	"context"
)

// RegisterRequest 商家管理员注册请求
type RegisterRequest struct {
	Email    string
	Mobile   string
	Password string
	RealName string
	ShopName string // 店铺名称，注册时创建
}

// LoginRequest 登录请求
type LoginRequest struct {
	Account  string // 邮箱、用户名或手机号
	Password string
	IP       string
}

// UpdateProfileRequest 更新个人资料
type UpdateProfileRequest struct {
	UserID   int64
	TenantID int64
	RealName string
	Avatar   string
	Mobile   string
	Email    string
}

// ChangePasswordRequest 修改密码
type ChangePasswordRequest struct {
	UserID          int64
	TenantID        int64
	OldPassword     string
	NewPassword     string
	ConfirmPassword string
}

// AdminUserResponse 管理员响应
type AdminUserResponse struct {
	ID          int64  `json:"id"`
	TenantID    int64  `json:"tenant_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	RealName    string `json:"real_name"`
	Avatar      string `json:"avatar"`
	Type        int    `json:"type"`
	TypeText    string `json:"type_text"`
	Status      int    `json:"status"`
	LastLoginAt string `json:"last_login_at"`
	CreatedAt   string `json:"created_at"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresIn    int64             `json:"expires_in"`
	User         AdminUserResponse `json:"user"`
}

// ListRequest 列表查询请求
type ListRequest struct {
	Page     int
	PageSize int
	TenantID int64  // 0表示查询所有（平台超管使用）
	Type     int    // 0表示不限制
	Status   int    // 0表示不限制
	Keyword  string // 搜索关键词
}

// ListResponse 列表响应
type ListResponse struct {
	List     []*AdminUserResponse `json:"list"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
PageSize int                  `json:"page_size"`
}

// CreateAdminUserRequest 创建管理员请求
type CreateAdminUserRequest struct {
	TenantID int64
Username string
Email    string
Password string
Mobile   string
RealName string
Avatar   string
Type     int
RoleIDs  []int64
}

// UpdateAdminUserRequest 更新管理员请求
type UpdateAdminUserRequest struct {
	ID       int64
RealName string
Avatar   string
Mobile   string
Email    string
}

// AssignRolesRequest 分配角色请求
type AssignRolesRequest struct {
AdminUserID int64
RoleIDs     []int64
}

// ResetPasswordResponse 重置密码响应
type ResetPasswordResponse struct {
TemporaryPassword string `json:"temporary_password"`
}

// Service 管理员应用服务接口
type Service interface {
	// RegisterTenantAdmin 商家管理员注册（自助开店）
	RegisterTenantAdmin(ctx context.Context, req RegisterRequest) (*LoginResponse, error)

	// Login 登录
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)

	// GetByID 根据ID获取管理员
	GetByID(ctx context.Context, id int64) (*AdminUserResponse, error)

	// GetByEmail 根据邮箱获取管理员
	GetByEmail(ctx context.Context, email string) (*AdminUserResponse, error)

	// UpdateProfile 更新个人资料
	UpdateProfile(ctx context.Context, req UpdateProfileRequest) (*AdminUserResponse, error)

	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, req ChangePasswordRequest) error

	// List 查询管理员列表（平台超管可查看所有，商家只能查看自己的）
	List(ctx context.Context, operatorID int64, req ListRequest) (*ListResponse, error)

	// Disable 禁用账号
	Disable(ctx context.Context, operatorID, targetID int64) error

	// Enable 启用账号
	Enable(ctx context.Context, operatorID, targetID int64) error

	// Create 创建管理员
	Create(ctx context.Context, operatorID int64, req CreateAdminUserRequest) (*AdminUserResponse, error)

	// Update 更新管理员
	Update(ctx context.Context, operatorID int64, req UpdateAdminUserRequest) (*AdminUserResponse, error)

	// Delete 删除管理员
	Delete(ctx context.Context, operatorID, targetID int64) error

	// AssignRoles 分配角色
	AssignRoles(ctx context.Context, operatorID int64, req AssignRolesRequest) error

	// ResetPassword 重置密码
	ResetPassword(ctx context.Context, operatorID, targetID int64) (*ResetPasswordResponse, error)
}
