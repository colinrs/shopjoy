package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/pkg/code"
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
	return db.WithContext(ctx).Where("deleted_at IS NULL").Save(u).Error
}

func (r *UserRepository) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	now := time.Now().UTC()
	dbQuery := db.WithContext(ctx).Model(&user.User{}).
		Where("id = ? AND deleted_at IS NULL", id)
	// Platform admin (tenantID == 0) can delete users across all tenants
	result := dbQuery.Update("deleted_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, db *gorm.DB, id int64) (*user.User, error) {
	var u user.User
	dbQuery := db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id)
	// Platform admin (tenantID == 0) can access users across all tenants
	err := dbQuery.First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.ErrUserNotFound
	}
	return &u, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, db *gorm.DB, email string) (*user.User, error) {
	var u user.User
	dbQuery := db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email)
	// Platform admin (tenantID == 0) can lookup by email across all tenants
	err := dbQuery.First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.ErrUserNotFound
	}
	return &u, err
}

func (r *UserRepository) FindByPhone(ctx context.Context, db *gorm.DB, phone string) (*user.User, error) {
	var u user.User
	dbQuery := db.WithContext(ctx).Where("phone = ? AND deleted_at IS NULL", phone)
	// Platform admin (tenantID == 0) can lookup by phone across all tenants
	err := dbQuery.First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.ErrUserNotFound
	}
	return &u, err
}

func (r *UserRepository) FindList(ctx context.Context, db *gorm.DB, query user.Query) ([]*user.User, int64, error) {
	var users []*user.User
	var total int64

	dbQuery := db.WithContext(ctx).Model(&user.User{}).Where("deleted_at IS NULL")

	// Platform admin (tenantID == 0) can access all tenants' data

	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Email != "" {
		dbQuery = dbQuery.Where("email = ?", query.Email)
	}
	if query.Phone != "" {
		dbQuery = dbQuery.Where("phone = ?", query.Phone)
	}
	if query.Status != 0 {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := dbQuery.Offset(query.Offset()).Limit(query.Limit()).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) Exists(ctx context.Context, db *gorm.DB, email, phone string) (bool, error) {
	var count int64
	dbQuery := db.WithContext(ctx).Model(&user.User{}).
		Where("deleted_at IS NULL AND (email = ? OR phone = ?)", email, phone)
	// Platform admin (tenantID == 0) can check duplicates across all tenants
	err := dbQuery.Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) GetStats(ctx context.Context, db *gorm.DB) (*user.Stats, error) {
	stats := &user.Stats{}

	// Total users (excluding deleted)
	if err := db.WithContext(ctx).Model(&user.User{}).
		Where("deleted_at IS NULL").
		Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Active users
	if err := db.WithContext(ctx).Model(&user.User{}).
		Where("status = ? AND deleted_at IS NULL", user.StatusActive).
		Count(&stats.Active).Error; err != nil {
		return nil, err
	}

	// Suspended users
	if err := db.WithContext(ctx).Model(&user.User{}).
		Where("status = ? AND deleted_at IS NULL", user.StatusSuspended).
		Count(&stats.Suspended).Error; err != nil {
		return nil, err
	}

	// New users today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if err := db.WithContext(ctx).Model(&user.User{}).
		Where("deleted_at IS NULL AND created_at >= ?", today).
		Count(&stats.NewToday).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
