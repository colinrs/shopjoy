package persistence

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() user.Repository {
	return &UserRepository{}
}

func (r *UserRepository) Create(ctx context.Context, db *gorm.DB, u *user.User) error {
	return db.WithContext(ctx).Create(u).Error
}

func (r *UserRepository) Update(ctx context.Context, db *gorm.DB, u *user.User) error {
	return db.WithContext(ctx).Save(u).Error
}

func (r *UserRepository) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*user.User, error) {
	var u user.User
	err := db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, user.ErrUserNotFound
	}
	return &u, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, email string) (*user.User, error) {
	var u user.User
	err := db.WithContext(ctx).Where("email = ? AND tenant_id = ?", email, tenantID.Int64()).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, user.ErrUserNotFound
	}
	return &u, err
}

func (r *UserRepository) FindByPhone(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, phone string) (*user.User, error) {
	var u user.User
	err := db.WithContext(ctx).Where("phone = ? AND tenant_id = ?", phone, tenantID.Int64()).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, user.ErrUserNotFound
	}
	return &u, err
}

func (r *UserRepository) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query user.Query) ([]*user.User, int64, error) {
	var users []*user.User
	var total int64

	dbQuery := db.WithContext(ctx).Model(&user.User{}).Where("tenant_id = ?", tenantID.Int64())

	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Email != "" {
		dbQuery = dbQuery.Where("email = ?", query.Email)
	}
	if query.Phone != "" {
		dbQuery = dbQuery.Where("phone = ?", query.Phone)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := dbQuery.Offset(query.Offset()).Limit(query.Limit()).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) Exists(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, email, phone string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(&user.User{}).
		Where("tenant_id = ? AND (email = ? OR phone = ?)", tenantID.Int64(), email, phone).
		Count(&count).Error
	return count > 0, err
}
