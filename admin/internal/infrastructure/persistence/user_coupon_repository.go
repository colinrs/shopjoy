package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type userCouponRepo struct{}

func NewUserCouponRepository() promotion.UserCouponRepository {
	return &userCouponRepo{}
}

// userCouponModel represents the database model for UserCoupon
type userCouponModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement:false"`
	TenantID   int64  `gorm:"column:tenant_id;not null;index"`
	UserID     int64  `gorm:"column:user_id;not null;index"`
	CouponID   int64  `gorm:"column:coupon_id;not null;index"`
	Status     int    `gorm:"column:status;not null;index"`
	UsedAt     *int64 `gorm:"column:used_at"`
	OrderID    string `gorm:"column:order_id;size:64"`
	ReceivedAt int64  `gorm:"column:received_at;not null"`
	ExpireAt   int64  `gorm:"column:expire_at;not null;index"`
	CreatedAt  int64  `gorm:"column:created_at"`
	UpdatedAt  int64  `gorm:"column:updated_at"`
}

func (userCouponModel) TableName() string {
	return "user_coupons"
}

func (m *userCouponModel) toEntity() *promotion.UserCoupon {
	var usedAt *time.Time
	if m.UsedAt != nil {
		t := time.Unix(*m.UsedAt, 0)
		usedAt = &t
	}

	return &promotion.UserCoupon{
		ID:         m.ID,
		TenantID:   shared.TenantID(m.TenantID),
		UserID:     m.UserID,
		CouponID:   m.CouponID,
		Status:     promotion.UserCouponStatus(m.Status),
		UsedAt:     usedAt,
		OrderID:    m.OrderID,
		ReceivedAt: time.Unix(m.ReceivedAt, 0),
		ExpireAt:   time.Unix(m.ExpireAt, 0),
		CreatedAt:  time.Unix(m.CreatedAt, 0),
		UpdatedAt:  time.Unix(m.UpdatedAt, 0),
	}
}

func fromUserCouponEntity(uc *promotion.UserCoupon) *userCouponModel {
	var usedAt *int64
	if uc.UsedAt != nil {
		ts := uc.UsedAt.Unix()
		usedAt = &ts
	}

	return &userCouponModel{
		ID:         uc.ID,
		TenantID:   uc.TenantID.Int64(),
		UserID:     uc.UserID,
		CouponID:   uc.CouponID,
		Status:     int(uc.Status),
		UsedAt:     usedAt,
		OrderID:    uc.OrderID,
		ReceivedAt: uc.ReceivedAt.Unix(),
		ExpireAt:   uc.ExpireAt.Unix(),
		CreatedAt:  uc.CreatedAt.Unix(),
		UpdatedAt:  uc.UpdatedAt.Unix(),
	}
}

// Create inserts a new user coupon
func (r *userCouponRepo) Create(ctx context.Context, db *gorm.DB, userCoupon *promotion.UserCoupon) error {
	model := fromUserCouponEntity(userCoupon)
	return db.WithContext(ctx).Create(model).Error
}

// FindByID finds a user coupon by ID
func (r *userCouponRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*promotion.UserCoupon, error) {
	query := db.WithContext(ctx).Where("tenant_id = ?", tenantID.Int64())
	var model userCouponModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrUserCouponNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindByUserID finds user coupons by user ID, optionally filtered by status
func (r *userCouponRepo) FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, status *promotion.UserCouponStatus) ([]*promotion.UserCoupon, error) {
	query := db.WithContext(ctx).Model(&userCouponModel{}).
		Where("tenant_id = ?", tenantID.Int64()).
		Where("user_id = ?", userID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var models []userCouponModel
	err := query.Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	userCoupons := make([]*promotion.UserCoupon, len(models))
	for i, m := range models {
		userCoupons[i] = m.toEntity()
	}
	return userCoupons, nil
}

// FindByUserAndCoupon finds user coupons by user ID and coupon ID
func (r *userCouponRepo) FindByUserAndCoupon(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, couponID int64) ([]*promotion.UserCoupon, error) {
	var models []userCouponModel
	err := db.WithContext(ctx).Model(&userCouponModel{}).
		Where("tenant_id = ?", tenantID.Int64()).
		Where("user_id = ?", userID).
		Where("coupon_id = ?", couponID).
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	userCoupons := make([]*promotion.UserCoupon, len(models))
	for i, m := range models {
		userCoupons[i] = m.toEntity()
	}
	return userCoupons, nil
}

// MarkUsed marks a user coupon as used
func (r *userCouponRepo) MarkUsed(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, orderID string) error {
	now := time.Now().UTC()
	result := db.WithContext(ctx).
		Model(&userCouponModel{}).
		Where("id = ? AND tenant_id = ? AND status = ?", id, tenantID.Int64(), promotion.UserCouponStatusUnused).
		Updates(map[string]interface{}{
			"status":     promotion.UserCouponStatusUsed,
			"used_at":    now,
			"order_id":   orderID,
			"updated_at": now,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrUserCouponNotFound
	}
	return nil
}

// CountUsageByUser counts how many times a user has used a specific coupon
func (r *userCouponRepo) CountUsageByUser(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, couponID int64) (int, error) {
	var count int64
	err := db.WithContext(ctx).Model(&userCouponModel{}).
		Where("tenant_id = ?", tenantID.Int64()).
		Where("user_id = ?", userID).
		Where("coupon_id = ?", couponID).
		Where("status = ?", promotion.UserCouponStatusUsed).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}