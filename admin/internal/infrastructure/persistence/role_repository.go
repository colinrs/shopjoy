package persistence

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/role"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
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

func (r *RoleRepository) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	return db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).Delete(&role.Role{}).Error
}

func (r *RoleRepository) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*role.Role, error) {
	var rl role.Role
	err := db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).First(&rl).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, role.ErrRoleNotFound
	}
	return &rl, err
}

func (r *RoleRepository) FindByCode(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, code string) (*role.Role, error) {
	var rl role.Role
	err := db.WithContext(ctx).Where("code = ? AND tenant_id = ?", code, tenantID.Int64()).First(&rl).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, role.ErrRoleNotFound
	}
	return &rl, err
}

func (r *RoleRepository) FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) ([]*role.Role, error) {
	var roles []*role.Role
	err := db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("roles.tenant_id = ? AND user_roles.user_id = ?", tenantID.Int64(), userID).
		Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query role.Query) ([]*role.Role, int64, error) {
	var roles []*role.Role
	var total int64

	dbQuery := db.WithContext(ctx).Model(&role.Role{}).Where("tenant_id = ?", tenantID.Int64())

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

func (r *RoleRepository) AssignToUser(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, roleIDs []int64) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&role.UserRole{}).Error; err != nil {
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

func (r *RoleRepository) GetUserRoles(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) ([]*role.Role, error) {
	return r.FindByUserID(ctx, db, tenantID, userID)
}
