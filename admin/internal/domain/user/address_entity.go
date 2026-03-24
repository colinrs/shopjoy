package user

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// UserAddress 用户收货地址实体
type UserAddress struct {
	ID         int64           `gorm:"column:id;primaryKey"`
	TenantID   shared.TenantID `gorm:"column:tenant_id;not null;index:idx_tenant_user"`
	UserID     int64           `gorm:"column:user_id;not null;index:idx_tenant_user"`
	Name       string          `gorm:"column:name;not null;size:100"`
	Phone      string          `gorm:"column:phone;not null;size:20"`
	Country    string          `gorm:"column:country;size:50;default:''"`
	Province   string          `gorm:"column:province;size:100;default:''"`
	City       string          `gorm:"column:city;size:100;default:''"`
	District   string          `gorm:"column:district;size:100;default:''"`
	Address    string          `gorm:"column:address;not null;size:255"`
	PostalCode string          `gorm:"column:postal_code;size:20;default:''"`
	IsDefault  bool            `gorm:"column:is_default;default:false"`
	DeletedAt  *int64          `gorm:"column:deleted_at"`
	CreatedAt  shared.UnixTime `gorm:"column:created_at;not null"`
	UpdatedAt  shared.UnixTime `gorm:"column:updated_at;not null"`
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