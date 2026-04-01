package adminuser

import (
	"context"

	"gorm.io/gorm"
)

// Repository 管理员仓储接口
type Repository interface {
	// Create 创建管理员
	Create(ctx context.Context, db *gorm.DB, user *AdminUser) error

	// Update 更新管理员
	Update(ctx context.Context, db *gorm.DB, user *AdminUser) error

	// Delete 软删除管理员
	Delete(ctx context.Context, db *gorm.DB, id int64) error

	// FindByID 根据ID查询
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*AdminUser, error)

	// FindByEmail 根据邮箱查询
	FindByEmail(ctx context.Context, db *gorm.DB, email string) (*AdminUser, error)

	// FindByUsername 根据用户名查询
	FindByUsername(ctx context.Context, db *gorm.DB, username string) (*AdminUser, error)

	// FindByMobile 根据手机号查询
	FindByMobile(ctx context.Context, db *gorm.DB, mobile string) (*AdminUser, error)

	// FindList 查询列表
	FindList(ctx context.Context, db *gorm.DB, query Query) ([]*AdminUser, int64, error)

	// Exists 检查邮箱或手机号是否已存在
	Exists(ctx context.Context, db *gorm.DB, email, mobile string) (bool, error)

	// UpdatePassword 更新密码
	UpdatePassword(ctx context.Context, db *gorm.DB, id int64, hashedPassword string) error

	// ExistsByUsername 检查用户名是否已存在
	ExistsByUsername(ctx context.Context, db *gorm.DB, tenantID int64, username string) (bool, error)

	// ExistsByEmail 检查邮箱是否已存在（在租户范围内）
	ExistsByEmail(ctx context.Context, db *gorm.DB, tenantID int64, email string) (bool, error)

	// CountMainAccount 统计租户主账号数量
	CountMainAccount(ctx context.Context, db *gorm.DB, tenantID int64) (int64, error)

	// CountByRoleID 统计拥有指定角色ID的用户数量
	CountByRoleID(ctx context.Context, db *gorm.DB, roleID int64) (int64, error)
}

// Query 查询条件
type Query struct {
	Page     int
	PageSize int
	TenantID int64  // 0表示查询所有（平台超管使用）
	Type     Type   // 0表示不限制
	Status   Status // 0表示不限制
	Keyword  string // 搜索关键词（用户名、邮箱、姓名）
}

func (q Query) Offset() int {
	return (q.Page - 1) * q.PageSize
}

func (q Query) Limit() int {
	if q.PageSize <= 0 {
		return 20
	}
	if q.PageSize > 100 {
		return 100
	}
	return q.PageSize
}
