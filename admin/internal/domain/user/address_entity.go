package user

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// UserAddress 用户收货地址实体
type UserAddress struct {
	ID         int64
	TenantID   shared.TenantID
	UserID     int64
	Name       string
	Phone      string
	Country    string
	Province   string
	City       string
	District   string
	Address    string
	PostalCode string
	IsDefault  bool
	DeletedAt  *int64
	CreatedAt  shared.UnixTime
	UpdatedAt  shared.UnixTime
}

func (a *UserAddress) TableName() string {
	return "user_addresses"
}

// SetAsDefault 设置为默认地址
func (a *UserAddress) SetAsDefault() {
	a.IsDefault = true
}

// UnsetDefault 取消默认地址
func (a *UserAddress) UnsetDefault() {
	a.IsDefault = false
}

// AddressQuery 地址查询参数
type AddressQuery struct {
	UserID int64
}

// AddressRepository 地址仓储接口
type AddressRepository interface {
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) ([]*UserAddress, error)
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*UserAddress, error)
}