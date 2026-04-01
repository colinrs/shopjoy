package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/shop"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ==================== Shop Settings Repository ====================

type shopSettingsRepo struct{}

func NewShopSettingsRepository() shop.ShopSettingsRepository {
	return &shopSettingsRepo{}
}

func (r *shopSettingsRepo) FindByTenantID(ctx context.Context, db *gorm.DB, tenantID int64) (*shop.ShopSettings, error) {
	var model shopSettingsModel
	err := db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *shopSettingsRepo) Save(ctx context.Context, db *gorm.DB, settings *shop.ShopSettings) error {
	// Check if exists
	var existing shopSettingsModel
	err := db.WithContext(ctx).
		Where("tenant_id = ?", settings.TenantID).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new
			now := time.Now().UTC()
			settings.CreatedAt = now
			settings.UpdatedAt = now
			model := fromShopSettingsEntity(settings)
			return db.WithContext(ctx).Create(model).Error
		}
		return err
	}

	// Update existing
	settings.ID = existing.ID
	settings.CreatedAt = existing.CreatedAt
	settings.UpdatedAt = time.Now().UTC()
	model := fromShopSettingsEntity(settings)
	return db.WithContext(ctx).Model(&shopSettingsModel{}).
		Where("id = ?", existing.ID).
		Updates(map[string]interface{}{
			"name":               model.Name,
			"code":               model.Code,
			"logo":               model.Logo,
			"description":        model.Description,
			"contact_name":       model.ContactName,
			"contact_phone":      model.ContactPhone,
			"contact_email":      model.ContactEmail,
			"address":            model.Address,
			"domain":             model.Domain,
			"custom_domain":       model.CustomDomain,
			"primary_color":      model.PrimaryColor,
			"secondary_color":    model.SecondaryColor,
			"favicon":            model.Favicon,
			"default_currency":    model.DefaultCurrency,
			"default_language":   model.DefaultLanguage,
			"timezone":           model.Timezone,
			"status":             model.Status,
			"plan":               model.Plan,
			"expire_at":          model.ExpireAt,
			"updated_at":         model.UpdatedAt,
		}).Error
}

type shopSettingsModel struct {
	ID               int64  `gorm:"column:id;primaryKey"`
	TenantID         int64  `gorm:"column:tenant_id;not null;uniqueIndex"`
	Name             string `gorm:"column:name;size:100;not null"`
	Code             string `gorm:"column:code;size:50;not null"`
	Logo             string `gorm:"column:logo;size:500"`
	Description      string `gorm:"column:description;size:500"`
	ContactName      string `gorm:"column:contact_name;size:50"`
	ContactPhone     string `gorm:"column:contact_phone;size:20"`
	ContactEmail     string `gorm:"column:contact_email;size:100"`
	Address          string `gorm:"column:address;size:255"`
	Domain           string `gorm:"column:domain;size:100"`
	CustomDomain     string `gorm:"column:custom_domain;size:100"`
	PrimaryColor     string `gorm:"column:primary_color;size:7"`
	SecondaryColor   string `gorm:"column:secondary_color;size:7"`
	Favicon          string `gorm:"column:favicon;size:255"`
	DefaultCurrency  string `gorm:"column:default_currency;size:3"`
	DefaultLanguage  string `gorm:"column:default_language;size:10"`
	Timezone         string `gorm:"column:timezone;size:50"`
	Status           int8   `gorm:"column:status;not null;default:1"`
	Plan             int8   `gorm:"column:plan;not null;default:0"`
	ExpireAt         string `gorm:"column:expire_at;size:50"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (m *shopSettingsModel) toEntity() *shop.ShopSettings {
	return &shop.ShopSettings{
		ID:               m.ID,
		TenantID:         m.TenantID,
		Name:             m.Name,
		Code:             m.Code,
		Logo:             m.Logo,
		Description:      m.Description,
		ContactName:      m.ContactName,
		ContactPhone:     m.ContactPhone,
		ContactEmail:     m.ContactEmail,
		Address:          m.Address,
		Domain:           m.Domain,
		CustomDomain:     m.CustomDomain,
		PrimaryColor:     m.PrimaryColor,
		SecondaryColor:   m.SecondaryColor,
		Favicon:          m.Favicon,
		DefaultCurrency:  m.DefaultCurrency,
		DefaultLanguage:  m.DefaultLanguage,
		Timezone:         m.Timezone,
		Status:           m.Status,
		Plan:             m.Plan,
		ExpireAt:         m.ExpireAt,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func fromShopSettingsEntity(s *shop.ShopSettings) *shopSettingsModel {
	return &shopSettingsModel{
		ID:               s.ID,
		TenantID:         s.TenantID,
		Name:             s.Name,
		Code:             s.Code,
		Logo:             s.Logo,
		Description:      s.Description,
		ContactName:      s.ContactName,
		ContactPhone:     s.ContactPhone,
		ContactEmail:     s.ContactEmail,
		Address:          s.Address,
		Domain:           s.Domain,
		CustomDomain:     s.CustomDomain,
		PrimaryColor:     s.PrimaryColor,
		SecondaryColor:   s.SecondaryColor,
		Favicon:          s.Favicon,
		DefaultCurrency:  s.DefaultCurrency,
		DefaultLanguage:  s.DefaultLanguage,
		Timezone:         s.Timezone,
		Status:           s.Status,
		Plan:             s.Plan,
		ExpireAt:         s.ExpireAt,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}

// ==================== Business Hours Repository ====================

type businessHoursRepo struct{}

func NewBusinessHoursRepository() shop.BusinessHoursRepository {
	return &businessHoursRepo{}
}

func (r *businessHoursRepo) FindByShopID(ctx context.Context, db *gorm.DB, shopID int64) ([]*shop.BusinessHours, error) {
	var models []businessHoursModel
	err := db.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Order("day_of_week ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*shop.BusinessHours, len(models))
	for i, m := range models {
		entities[i] = m.toEntity()
	}
	return entities, nil
}

func (r *businessHoursRepo) SaveBatch(ctx context.Context, db *gorm.DB, shopID int64, hours []*shop.BusinessHours) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing hours for this shop
		if err := tx.Where("shop_id = ?", shopID).Delete(&businessHoursModel{}).Error; err != nil {
			return err
		}

		// Insert new hours
		now := time.Now().UTC()
		for _, h := range hours {
			h.ShopID = shopID
			h.CreatedAt = now
			h.UpdatedAt = now
			model := fromBusinessHoursEntity(h)
			if err := tx.Create(model).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

type businessHoursModel struct {
	ID        int64  `gorm:"column:id;primaryKey"`
	ShopID    int64  `gorm:"column:shop_id;not null;index"`
	DayOfWeek int8   `gorm:"column:day_of_week;not null"`
	OpenTime  string `gorm:"column:open_time;size:5"`
	CloseTime string `gorm:"column:close_time;size:5"`
	IsClosed  bool   `gorm:"column:is_closed;not null;default:false"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (m *businessHoursModel) toEntity() *shop.BusinessHours {
	return &shop.BusinessHours{
		ID:        m.ID,
		ShopID:    m.ShopID,
		DayOfWeek: m.DayOfWeek,
		OpenTime:  m.OpenTime,
		CloseTime: m.CloseTime,
		IsClosed:  m.IsClosed,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func fromBusinessHoursEntity(h *shop.BusinessHours) *businessHoursModel {
	return &businessHoursModel{
		ID:        h.ID,
		ShopID:    h.ShopID,
		DayOfWeek: h.DayOfWeek,
		OpenTime:  h.OpenTime,
		CloseTime: h.CloseTime,
		IsClosed:  h.IsClosed,
		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
	}
}

// ==================== Notification Settings Repository ====================

type notificationSettingsRepo struct{}

func NewNotificationSettingsRepository() shop.NotificationSettingsRepository {
	return &notificationSettingsRepo{}
}

func (r *notificationSettingsRepo) FindByShopID(ctx context.Context, db *gorm.DB, shopID int64) (*shop.NotificationSettings, error) {
	var model notificationSettingsModel
	err := db.WithContext(ctx).
		Where("shop_id = ?", shopID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *notificationSettingsRepo) Save(ctx context.Context, db *gorm.DB, settings *shop.NotificationSettings) error {
	// Check if exists
	var existing notificationSettingsModel
	err := db.WithContext(ctx).
		Where("shop_id = ?", settings.ShopID).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new
			now := time.Now().UTC()
			settings.CreatedAt = now
			settings.UpdatedAt = now
			model := fromNotificationSettingsEntity(settings)
			return db.WithContext(ctx).Create(model).Error
		}
		return err
	}

	// Update existing
	settings.ID = existing.ID
	settings.CreatedAt = existing.CreatedAt
	settings.UpdatedAt = time.Now().UTC()
	model := fromNotificationSettingsEntity(settings)
	return db.WithContext(ctx).Model(&notificationSettingsModel{}).
		Where("id = ?", existing.ID).
		Updates(map[string]interface{}{
			"order_created":       model.OrderCreated,
			"order_paid":          model.OrderPaid,
			"order_shipped":       model.OrderShipped,
			"order_cancelled":     model.OrderCancelled,
			"low_stock_alert":     model.LowStockAlert,
			"low_stock_threshold": model.LowStockThreshold,
			"refund_requested":    model.RefundRequested,
			"new_review":          model.NewReview,
			"updated_at":          model.UpdatedAt,
		}).Error
}

type notificationSettingsModel struct {
	ID                int64 `gorm:"column:id;primaryKey"`
	ShopID            int64 `gorm:"column:shop_id;not null;uniqueIndex"`
	OrderCreated      bool  `gorm:"column:order_created;not null;default:true"`
	OrderPaid         bool  `gorm:"column:order_paid;not null;default:true"`
	OrderShipped      bool  `gorm:"column:order_shipped;not null;default:true"`
	OrderCancelled    bool  `gorm:"column:order_cancelled;not null;default:true"`
	LowStockAlert     bool  `gorm:"column:low_stock_alert;not null;default:true"`
	LowStockThreshold int   `gorm:"column:low_stock_threshold;not null;default:10"`
	RefundRequested   bool  `gorm:"column:refund_requested;not null;default:true"`
	NewReview         bool  `gorm:"column:new_review;not null;default:true"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

func (m *notificationSettingsModel) toEntity() *shop.NotificationSettings {
	return &shop.NotificationSettings{
		ID:                m.ID,
		ShopID:            m.ShopID,
		OrderCreated:      m.OrderCreated,
		OrderPaid:         m.OrderPaid,
		OrderShipped:      m.OrderShipped,
		OrderCancelled:    m.OrderCancelled,
		LowStockAlert:     m.LowStockAlert,
		LowStockThreshold: m.LowStockThreshold,
		RefundRequested:   m.RefundRequested,
		NewReview:         m.NewReview,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func fromNotificationSettingsEntity(s *shop.NotificationSettings) *notificationSettingsModel {
	return &notificationSettingsModel{
		ID:                s.ID,
		ShopID:            s.ShopID,
		OrderCreated:      s.OrderCreated,
		OrderPaid:         s.OrderPaid,
		OrderShipped:      s.OrderShipped,
		OrderCancelled:    s.OrderCancelled,
		LowStockAlert:     s.LowStockAlert,
		LowStockThreshold: s.LowStockThreshold,
		RefundRequested:   s.RefundRequested,
		NewReview:         s.NewReview,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

// ==================== Payment Settings Repository ====================

type paymentSettingsRepo struct{}

func NewPaymentSettingsRepository() shop.PaymentSettingsRepository {
	return &paymentSettingsRepo{}
}

func (r *paymentSettingsRepo) FindByShopID(ctx context.Context, db *gorm.DB, shopID int64) (*shop.PaymentSettings, error) {
	var model paymentSettingsModel
	err := db.WithContext(ctx).
		Where("shop_id = ?", shopID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *paymentSettingsRepo) Save(ctx context.Context, db *gorm.DB, settings *shop.PaymentSettings) error {
	// Check if exists
	var existing paymentSettingsModel
	err := db.WithContext(ctx).
		Where("shop_id = ?", settings.ShopID).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new
			now := time.Now().UTC()
			settings.CreatedAt = now
			settings.UpdatedAt = now
			model := fromPaymentSettingsEntity(settings)
			return db.WithContext(ctx).Create(model).Error
		}
		return err
	}

	// Update existing
	settings.ID = existing.ID
	settings.CreatedAt = existing.CreatedAt
	settings.UpdatedAt = time.Now().UTC()
	model := fromPaymentSettingsEntity(settings)
	return db.WithContext(ctx).Model(&paymentSettingsModel{}).
		Where("id = ?", existing.ID).
		Updates(map[string]interface{}{
			"stripe_enabled":    model.StripeEnabled,
			"stripe_public_key": model.StripePublicKey,
			"stripe_secret_key": model.StripeSecretKey,
			"updated_at":        model.UpdatedAt,
		}).Error
}

type paymentSettingsModel struct {
	ID              int64     `gorm:"column:id;primaryKey"`
	ShopID          int64     `gorm:"column:shop_id;not null;uniqueIndex"`
	StripeEnabled   bool      `gorm:"column:stripe_enabled;not null;default:false"`
	StripePublicKey string    `gorm:"column:stripe_public_key;size:255"`
	StripeSecretKey string    `gorm:"column:stripe_secret_key;size:255"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (m *paymentSettingsModel) toEntity() *shop.PaymentSettings {
	return &shop.PaymentSettings{
		ID:              m.ID,
		ShopID:          m.ShopID,
		StripeEnabled:   m.StripeEnabled,
		StripePublicKey: m.StripePublicKey,
		StripeSecretKey: m.StripeSecretKey,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func fromPaymentSettingsEntity(s *shop.PaymentSettings) *paymentSettingsModel {
	return &paymentSettingsModel{
		ID:              s.ID,
		ShopID:          s.ShopID,
		StripeEnabled:   s.StripeEnabled,
		StripePublicKey: s.StripePublicKey,
		StripeSecretKey: s.StripeSecretKey,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}

// ==================== Shipping Settings Repository ====================

type shippingSettingsRepo struct{}

func NewShippingSettingsRepository() shop.ShippingSettingsRepository {
	return &shippingSettingsRepo{}
}

func (r *shippingSettingsRepo) FindByShopID(ctx context.Context, db *gorm.DB, shopID int64) (*shop.ShippingSettings, error) {
	var model shippingSettingsModel
	err := db.WithContext(ctx).
		Where("shop_id = ?", shopID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *shippingSettingsRepo) Save(ctx context.Context, db *gorm.DB, settings *shop.ShippingSettings) error {
	// Check if exists
	var existing shippingSettingsModel
	err := db.WithContext(ctx).
		Where("shop_id = ?", settings.ShopID).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new
			now := time.Now().UTC()
			settings.CreatedAt = now
			settings.UpdatedAt = now
			model := fromShippingSettingsEntity(settings)
			return db.WithContext(ctx).Create(model).Error
		}
		return err
	}

	// Update existing
	settings.ID = existing.ID
	settings.CreatedAt = existing.CreatedAt
	settings.UpdatedAt = time.Now().UTC()
	model := fromShippingSettingsEntity(settings)
	return db.WithContext(ctx).Model(&shippingSettingsModel{}).
		Where("id = ?", existing.ID).
		Updates(map[string]interface{}{
			"free_shipping_threshold": model.FreeShippingThreshold,
			"default_shipping_fee":    model.DefaultShippingFee,
			"currency":                model.Currency,
			"updated_at":             model.UpdatedAt,
		}).Error
}

type shippingSettingsModel struct {
	ID                     int64     `gorm:"column:id;primaryKey"`
	ShopID                 int64     `gorm:"column:shop_id;not null;uniqueIndex"`
	FreeShippingThreshold  string    `gorm:"column:free_shipping_threshold;type:decimal(19,4);not null"`
	DefaultShippingFee     string    `gorm:"column:default_shipping_fee;type:decimal(19,4);not null"`
	Currency               string    `gorm:"column:currency;size:3;not null"`
	CreatedAt              time.Time `gorm:"column:created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at"`
}

func (m *shippingSettingsModel) toEntity() *shop.ShippingSettings {
	threshold, _ := decimal.NewFromString(m.FreeShippingThreshold)
	fee, _ := decimal.NewFromString(m.DefaultShippingFee)
	return &shop.ShippingSettings{
		ID:                     m.ID,
		ShopID:                 m.ShopID,
		FreeShippingThreshold:  threshold,
		DefaultShippingFee:     fee,
		Currency:               m.Currency,
		CreatedAt:              m.CreatedAt,
		UpdatedAt:              m.UpdatedAt,
	}
}

func fromShippingSettingsEntity(s *shop.ShippingSettings) *shippingSettingsModel {
	return &shippingSettingsModel{
		ID:                     s.ID,
		ShopID:                 s.ShopID,
		FreeShippingThreshold:  s.FreeShippingThreshold.String(),
		DefaultShippingFee:     s.DefaultShippingFee.String(),
		Currency:               s.Currency,
		CreatedAt:              s.CreatedAt,
		UpdatedAt:              s.UpdatedAt,
	}
}
