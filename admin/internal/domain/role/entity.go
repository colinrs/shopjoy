package role

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type Status int

const (
	StatusDisabled Status = iota
	StatusEnabled
)

type Role struct {
	application.Model
	TenantID    shared.TenantID
	Name        string
	Code        string
	Description string
	Status      Status
	IsSystem    bool
	Audit       shared.AuditInfo `gorm:"embedded"`
}

func (r *Role) TableName() string {
	return "roles"
}

func (r *Role) Enable() {
	r.Status = StatusEnabled
}

func (r *Role) Disable() {
	r.Status = StatusDisabled
}

func (r *Role) IsActive() bool {
	return r.Status == StatusEnabled
}

type Permission struct {
	application.Model
	Name     string
	Code     string
	Type     PermissionType
	ParentID int64
	Path     string
	Icon     string
	Sort     int
}

type PermissionType int

const (
	PermissionTypeMenu PermissionType = iota
	PermissionTypeButton
	PermissionTypeAPI
)

type RolePermission struct {
	application.Model
	RoleID       int64
	PermissionID int64
}

type UserRole struct {
	application.Model
	UserID int64
	RoleID int64
}

type Repository interface {
	Create(ctx context.Context, db *gorm.DB, role *Role) error
	Update(ctx context.Context, db *gorm.DB, role *Role) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Role, error)
	FindByCode(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, code string) (*Role, error)
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) ([]*Role, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query Query) ([]*Role, int64, error)
	AssignToUser(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, roleIDs []int64) error
	GetUserRoles(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) ([]*Role, error)
}

type PermissionRepository interface {
	FindAll(ctx context.Context, db *gorm.DB) ([]*Permission, error)
	FindByRoleIDs(ctx context.Context, db *gorm.DB, roleIDs []int64) ([]*Permission, error)
	AssignToRole(ctx context.Context, db *gorm.DB, roleID int64, permissionIDs []int64) error
}

type Query struct {
	shared.PageQuery
	Name   string
	Code   string
	Status Status
}
