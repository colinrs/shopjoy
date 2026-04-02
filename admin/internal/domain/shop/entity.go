package shop

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ShopSettings 店铺设置
type ShopSettings struct {
	ID              int64     `gorm:"column:id;primaryKey"`
	TenantID        int64     `gorm:"column:tenant_id;not null;uniqueIndex"`
	Name            string    `gorm:"column:name;size:100;not null"`
	Code            string    `gorm:"column:code;size:50;not null"`
	Logo            string    `gorm:"column:logo;size:500"`
	Description     string    `gorm:"column:description;size:500"`
	ContactName     string    `gorm:"column:contact_name;size:50"`
	ContactPhone    string    `gorm:"column:contact_phone;size:20"`
	ContactEmail    string    `gorm:"column:contact_email;size:100"`
	Address         string    `gorm:"column:address;size:255"`
	Domain          string    `gorm:"column:domain;size:100"`
	CustomDomain    string    `gorm:"column:custom_domain;size:100"`
	PrimaryColor    string    `gorm:"column:primary_color;size:7"`
	SecondaryColor  string    `gorm:"column:secondary_color;size:7"`
	Favicon         string    `gorm:"column:favicon;size:255"`
	DefaultCurrency string    `gorm:"column:default_currency;size:3"`
	DefaultLanguage string    `gorm:"column:default_language;size:10"`
	Timezone        string    `gorm:"column:timezone;size:50"`
	Status          int8      `gorm:"column:status;not null;default:1"`
	Plan            int8      `gorm:"column:plan;not null;default:0"`
	ExpireAt        string    `gorm:"column:expire_at;size:50"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (s *ShopSettings) TableName() string {
	return "shop_settings"
}

// StatusText returns the status text representation
func (s *ShopSettings) StatusText() string {
	switch s.Status {
	case 0:
		return "inactive"
	case 1:
		return "active"
	case 2:
		return "suspended"
	default:
		return "unknown"
	}
}

// PlanText returns the plan text representation
func (s *ShopSettings) PlanText() string {
	switch s.Plan {
	case 0:
		return "basic"
	case 1:
		return "standard"
	case 2:
		return "premium"
	default:
		return "unknown"
	}
}

// BusinessHours 营业时间
type BusinessHours struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	ShopID    int64     `gorm:"column:shop_id;not null;index"`
	DayOfWeek int8      `gorm:"column:day_of_week;not null"`
	OpenTime  string    `gorm:"column:open_time;size:5"`
	CloseTime string    `gorm:"column:close_time;size:5"`
	IsClosed  bool      `gorm:"column:is_closed;not null;default:false"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (b *BusinessHours) TableName() string {
	return "shop_business_hours"
}

// NotificationSettings 通知设置
type NotificationSettings struct {
	ID                int64     `gorm:"column:id;primaryKey"`
	ShopID            int64     `gorm:"column:shop_id;not null;uniqueIndex"`
	OrderCreated      bool      `gorm:"column:order_created;not null;default:true"`
	OrderPaid         bool      `gorm:"column:order_paid;not null;default:true"`
	OrderShipped      bool      `gorm:"column:order_shipped;not null;default:true"`
	OrderCancelled    bool      `gorm:"column:order_cancelled;not null;default:true"`
	LowStockAlert     bool      `gorm:"column:low_stock_alert;not null;default:true"`
	LowStockThreshold int       `gorm:"column:low_stock_threshold;not null;default:10"`
	RefundRequested   bool      `gorm:"column:refund_requested;not null;default:true"`
	NewReview         bool      `gorm:"column:new_review;not null;default:true"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

func (n *NotificationSettings) TableName() string {
	return "shop_notification_settings"
}

// PaymentSettings 支付设置
type PaymentSettings struct {
	ID              int64     `gorm:"column:id;primaryKey"`
	ShopID          int64     `gorm:"column:shop_id;not null;uniqueIndex"`
	StripeEnabled   bool      `gorm:"column:stripe_enabled;not null;default:false"`
	StripePublicKey string    `gorm:"column:stripe_public_key;size:255"`
	StripeSecretKey string    `gorm:"column:stripe_secret_key;size:255"` // Encrypted
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (p *PaymentSettings) TableName() string {
	return "shop_payment_settings"
}

// ShippingSettings 运费设置
type ShippingSettings struct {
	ID                    int64           `gorm:"column:id;primaryKey"`
	ShopID                int64           `gorm:"column:shop_id;not null;uniqueIndex"`
	FreeShippingThreshold decimal.Decimal `gorm:"column:free_shipping_threshold;type:decimal(19,4);not null"`
	DefaultShippingFee    decimal.Decimal `gorm:"column:default_shipping_fee;type:decimal(19,4);not null"`
	Currency              string          `gorm:"column:currency;size:3;not null"`
	CreatedAt             time.Time       `gorm:"column:created_at"`
	UpdatedAt             time.Time       `gorm:"column:updated_at"`
}

func (s *ShippingSettings) TableName() string {
	return "shop_shipping_settings"
}

// Repository interfaces

type ShopSettingsRepository interface {
	FindByTenantID(ctx context.Context, db *gorm.DB, tenantID int64) (*ShopSettings, error)
	Save(ctx context.Context, db *gorm.DB, settings *ShopSettings) error
}

type BusinessHoursRepository interface {
	FindByShopID(ctx context.Context, db *gorm.DB, shopID int64) ([]*BusinessHours, error)
	SaveBatch(ctx context.Context, db *gorm.DB, shopID int64, hours []*BusinessHours) error
}

type NotificationSettingsRepository interface {
	FindByShopID(ctx context.Context, db *gorm.DB, shopID int64) (*NotificationSettings, error)
	Save(ctx context.Context, db *gorm.DB, settings *NotificationSettings) error
}

type PaymentSettingsRepository interface {
	FindByShopID(ctx context.Context, db *gorm.DB, shopID int64) (*PaymentSettings, error)
	Save(ctx context.Context, db *gorm.DB, settings *PaymentSettings) error
}

type ShippingSettingsRepository interface {
	FindByShopID(ctx context.Context, db *gorm.DB, shopID int64) (*ShippingSettings, error)
	Save(ctx context.Context, db *gorm.DB, settings *ShippingSettings) error
}
