package persistence

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/role"
	"github.com/colinrs/shopjoy/pkg/code"
	"gorm.io/gorm"
)

type RoleRepository struct{}

func NewRoleRepository() role.Repository {
	return &RoleRepository{}
}

func (r *RoleRepository) Create(ctx context.Context, db *gorm.DB, role *role.Role) error {
	return db.WithContext(ctx).Create(role).Error
}

func (r *RoleRepository) Update(ctx context.Context, db *gorm.DB, role *role.Role) error {
	return db.WithContext(ctx).Save(role).Error
}

func (r *RoleRepository) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	return db.WithContext(ctx).Where("id = ?", id).Delete(&role.Role{}).Error
}

func (r *RoleRepository) FindByID(ctx context.Context, db *gorm.DB, id int64) (*role.Role, error) {
	var rl role.Role
	err := db.WithContext(ctx).Where("id = ?", id).First(&rl).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.ErrRoleNotFound
	}
	return &rl, err
}

func (r *RoleRepository) FindByCode(ctx context.Context, db *gorm.DB, codeStr string) (*role.Role, error) {
	var rl role.Role
	err := db.WithContext(ctx).Where("code = ?", codeStr).First(&rl).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.ErrRoleNotFound
	}
	return &rl, err
}

func (r *RoleRepository) FindByUserID(ctx context.Context, db *gorm.DB, userID int64) ([]*role.Role, error) {
	var roles []*role.Role
	err := db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) FindList(ctx context.Context, db *gorm.DB, query role.Query) ([]*role.Role, int64, error) {
	var roles []*role.Role
	var total int64

	dbQuery := db.WithContext(ctx).Model(&role.Role{})

	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Code != "" {
		dbQuery = dbQuery.Where("code = ?", query.Code)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := dbQuery.Offset(query.Offset()).Limit(query.Limit()).Find(&roles).Error
	return roles, total, err
}

func (r *RoleRepository) AssignToUser(ctx context.Context, db *gorm.DB, userID int64, roleIDs []int64) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Junction table user_roles has PRIMARY KEY (user_id, role_id) which
		// excludes deleted_at. A soft delete (default for application.Model)
		// would only set deleted_at, leaving the physical row in place and
		// causing the next INSERT to collide on the PK (Error 1062). Use
		// Unscoped() to physically remove prior rows before re-inserting.
		if err := tx.Where("user_id = ?", userID).Unscoped().Delete(&role.UserRole{}).Error; err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			userRole := &role.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
			if err := tx.Create(userRole).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RoleRepository) GetUserRoles(ctx context.Context, db *gorm.DB, userID int64) ([]*role.Role, error) {
	return r.FindByUserID(ctx, db, userID)
}

// PermissionRepository implements role.PermissionRepository
// Note: Permissions are global (not tenant-scoped) and shared across all tenants.
// This design allows for centralized permission management while roles are tenant-specific.
type PermissionRepository struct{}

func NewPermissionRepository() role.PermissionRepository {
	return &PermissionRepository{}
}

func (r *PermissionRepository) FindAll(ctx context.Context, db *gorm.DB) ([]*role.Permission, error) {
	var permissions []*role.Permission
	err := db.WithContext(ctx).Order("sort ASC").Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) FindByRoleIDs(ctx context.Context, db *gorm.DB, roleIDs []int64) ([]*role.Permission, error) {
	if len(roleIDs) == 0 {
		return []*role.Permission{}, nil
	}
	var permissions []*role.Permission
	err := db.WithContext(ctx).
		Distinct().
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Order("permissions.sort ASC").
		Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) AssignToRole(ctx context.Context, db *gorm.DB, roleID int64, permissionIDs []int64) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Junction table role_permissions has PRIMARY KEY (role_id, permission_id)
		// which excludes deleted_at; soft-delete would leave the row in place and
		// collide on the PK on the next INSERT. See AssignToUser for full rationale.
		if err := tx.Where("role_id = ?", roleID).Unscoped().Delete(&role.RolePermission{}).Error; err != nil {
			return err
		}
		for _, permID := range permissionIDs {
			rp := &role.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
			}
			if err := tx.Create(rp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
