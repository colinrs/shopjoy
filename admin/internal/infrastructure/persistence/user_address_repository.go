package persistence

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/pkg/code"
	"gorm.io/gorm"
)

type UserAddressRepository struct{}

func NewUserAddressRepository() user.AddressRepository {
	return &UserAddressRepository{}
}

func (r *UserAddressRepository) FindByUserID(ctx context.Context, db *gorm.DB,  userID int64) ([]*user.UserAddress, error) {
	var addresses []*user.UserAddress
	err := db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *UserAddressRepository) FindByID(ctx context.Context, db *gorm.DB,  id int64) (*user.UserAddress, error) {
	var address user.UserAddress
	err := db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&address).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.ErrAddressNotFound
	}
	if err != nil {
		return nil, err
	}
	return &address, nil
}
