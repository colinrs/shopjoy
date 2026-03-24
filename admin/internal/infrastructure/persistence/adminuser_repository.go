package persistence

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/adminuser"
	"github.com/colinrs/shopjoy/pkg/code"
	"gorm.io/gorm"
)

type AdminUserRepository struct{}

func NewAdminUserRepository() adminuser.Repository {
	return &AdminUserRepository{}
}

func (r *AdminUserRepository) Create(ctx context.Context, db *gorm.DB, user *adminuser.AdminUser) error {
	return db.WithContext(ctx).Create(user).Error
}

func (r *AdminUserRepository) Update(ctx context.Context, db *gorm.DB, user *adminuser.AdminUser) error {
	return db.WithContext(ctx).Save(user).Error
}

func (r *AdminUserRepository) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	return db.WithContext(ctx).Delete(&adminuser.AdminUser{}, id).Error
}

func (r *AdminUserRepository) FindByID(ctx context.Context, db *gorm.DB, id int64) (*adminuser.AdminUser, error) {
	var user adminuser.AdminUser
	err := db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.ErrAdminUserNotFound
	}
	return &user, err
}

func (r *AdminUserRepository) FindByEmail(ctx context.Context, db *gorm.DB, email string) (*adminuser.AdminUser, error) {
	var user adminuser.AdminUser
	err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.ErrAdminUserNotFound
	}
	return &user, err
}

func (r *AdminUserRepository) FindByUsername(ctx context.Context, db *gorm.DB, username string) (*adminuser.AdminUser, error) {
	var user adminuser.AdminUser
	err := db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.ErrAdminUserNotFound
	}
	return &user, err
}

func (r *AdminUserRepository) FindByMobile(ctx context.Context, db *gorm.DB, mobile string) (*adminuser.AdminUser, error) {
	var user adminuser.AdminUser
	err := db.WithContext(ctx).Where("mobile = ?", mobile).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.ErrAdminUserNotFound
	}
	return &user, err
}

func (r *AdminUserRepository) FindList(ctx context.Context, db *gorm.DB, query adminuser.Query) ([]*adminuser.AdminUser, int64, error) {
	var users []*adminuser.AdminUser
	var total int64

	dbQuery := db.WithContext(ctx).Model(&adminuser.AdminUser{})

	if query.TenantID > 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", query.TenantID)
	}

	if query.Type > 0 {
		dbQuery = dbQuery.Where("type = ?", query.Type)
	}

	if query.Status > 0 {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		dbQuery = dbQuery.Where("username LIKE ? OR email LIKE ? OR real_name LIKE ?", keyword, keyword, keyword)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := dbQuery.Offset(query.Offset()).Limit(query.Limit()).Order("id DESC").Find(&users).Error
	return users, total, err
}

func (r *AdminUserRepository) Exists(ctx context.Context, db *gorm.DB, email, mobile string) (bool, error) {
	var count int64
	dbQuery := db.WithContext(ctx).Model(&adminuser.AdminUser{})

	if email != "" && mobile != "" {
		dbQuery = dbQuery.Where("email = ? OR mobile = ?", email, mobile)
	} else if email != "" {
		dbQuery = dbQuery.Where("email = ?", email)
	} else if mobile != "" {
		dbQuery = dbQuery.Where("mobile = ?", mobile)
	}

	err := dbQuery.Count(&count).Error
	return count > 0, err
}

func (r *AdminUserRepository) UpdatePassword(ctx context.Context, db *gorm.DB, id int64, hashedPassword string) error {
	return db.WithContext(ctx).Model(&adminuser.AdminUser{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error
}

func (r *AdminUserRepository) ExistsByUsername(ctx context.Context, db *gorm.DB, tenantID int64, username string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(&adminuser.AdminUser{}).
		Where("tenant_id = ? AND username = ?", tenantID, username).
		Count(&count).Error
	return count > 0, err
}

func (r *AdminUserRepository) ExistsByEmail(ctx context.Context, db *gorm.DB, tenantID int64, email string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(&adminuser.AdminUser{}).
		Where("tenant_id = ? AND email = ?", tenantID, email).
		Count(&count).Error
	return count > 0, err
}

func (r *AdminUserRepository) CountMainAccount(ctx context.Context, db *gorm.DB, tenantID int64) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Model(&adminuser.AdminUser{}).
		Where("tenant_id = ? AND type = ?", tenantID, adminuser.TypeTenantAdmin).
		Count(&count).Error
	return count, err
}
