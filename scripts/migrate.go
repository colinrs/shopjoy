package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/colinrs/shopjoy/pkg/infra"
	"github.com/zeromicro/go-zero/core/conf"
	"gorm.io/gorm"
)

var configFile = flag.String("f", "etc/admin-api.yaml", "配置文件路径")

type Config struct {
	MySQL infra.DBConfig
}

func main1() {
	flag.Parse()

	var c Config
	conf.MustLoad(*configFile, &c)

	db, err := infra.Database(&c.MySQL)
	if err != nil {
		fmt.Printf("连接数据库失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("开始自动迁移数据库...")

	if err := migrate(db); err != nil {
		fmt.Printf("迁移失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("数据库迁移完成!")
}

func migrate(db *gorm.DB) error {
	models := []interface{}{
		&Tenant{},
		&User{},
		&Role{},
		&UserRole{},
		&Product{},
		&SKU{},
		&Category{},
		&Brand{},
		&Order{},
		&OrderItem{},
		&Cart{},
		&CartItem{},
		&Payment{},
		&Coupon{},
		&UserCoupon{},
		&Shop{},
		&Theme{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("migrate %T failed: %w", model, err)
		}
	}

	return nil
}

type Tenant struct {
	ID           int64  `gorm:"primaryKey"`
	Name         string `gorm:"size:255;not null"`
	Code         string `gorm:"size:100;uniqueIndex;not null"`
	Status       int8
	Plan         int8
	Domain       string `gorm:"size:255"`
	CustomDomain string `gorm:"size:255"`
	Logo         string `gorm:"size:500"`
	ContactName  string `gorm:"size:100"`
	ContactPhone string `gorm:"size:20"`
	ContactEmail string `gorm:"size:255"`
	Address      string `gorm:"type:text"`
	ExpireAt     *int64
	CreatedAt    int64 `gorm:"autoCreateTime"`
	UpdatedAt    int64 `gorm:"autoUpdateTime"`
}

func (Tenant) TableName() string {
	return "tenants"
}

type User struct {
	ID        int64  `gorm:"primaryKey"`
	TenantID  int64  `gorm:"index;not null"`
	Email     string `gorm:"size:255;not null"`
	Phone     string `gorm:"size:20"`
	Password  string `gorm:"size:255;not null"`
	Name      string `gorm:"size:100;not null"`
	Avatar    string `gorm:"size:500"`
	Gender    int8
	Birthday  *int64
	Status    int8
	LastLogin *int64
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
	CreatedBy int64
	UpdatedBy int64
}

func (User) TableName() string {
	return "users"
}

type Role struct {
	ID          int64  `gorm:"primaryKey"`
	TenantID    int64  `gorm:"index;not null"`
	Name        string `gorm:"size:100;not null"`
	Code        string `gorm:"size:100;not null"`
	Description string `gorm:"type:text"`
	Status      int8
	IsSystem    int8
	CreatedAt   int64 `gorm:"autoCreateTime"`
	UpdatedAt   int64 `gorm:"autoUpdateTime"`
	CreatedBy   int64
	UpdatedBy   int64
}

func (Role) TableName() string {
	return "roles"
}

type UserRole struct {
	UserID int64 `gorm:"primaryKey"`
	RoleID int64 `gorm:"primaryKey"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

type Product struct {
	ID          int64  `gorm:"primaryKey"`
	TenantID    int64  `gorm:"index;not null"`
	Name        string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	Price       int64
	CostPrice   int64
	Stock       int
	Status      int8
	CategoryID  int64  `gorm:"index"`
	BrandID     int64  `gorm:"index"`
	Images      string `gorm:"type:json"`
	Attributes  string `gorm:"type:json"`
	CreatedAt   int64  `gorm:"autoCreateTime"`
	UpdatedAt   int64  `gorm:"autoUpdateTime"`
}

func (Product) TableName() string {
	return "products"
}

type SKU struct {
	ID         int64  `gorm:"primaryKey"`
	ProductID  int64  `gorm:"index;not null"`
	Code       string `gorm:"size:100;uniqueIndex;not null"`
	Price      int64
	Stock      int
	Attributes string `gorm:"type:json"`
	Status     int8
	CreatedAt  int64 `gorm:"autoCreateTime"`
	UpdatedAt  int64 `gorm:"autoUpdateTime"`
}

func (SKU) TableName() string {
	return "skus"
}

type Category struct {
	ID        int64  `gorm:"primaryKey"`
	TenantID  int64  `gorm:"index;not null"`
	ParentID  int64  `gorm:"index;default:0"`
	Name      string `gorm:"size:100;not null"`
	Code      string `gorm:"size:100"`
	Level     int8
	Sort      int
	Icon      string `gorm:"size:255"`
	Image     string `gorm:"size:500"`
	Status    int8
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
}

func (Category) TableName() string {
	return "categories"
}

type Brand struct {
	ID          int64  `gorm:"primaryKey"`
	TenantID    int64  `gorm:"index;not null"`
	Name        string `gorm:"size:100;not null"`
	Logo        string `gorm:"size:500"`
	Description string `gorm:"type:text"`
	Website     string `gorm:"size:255"`
	Sort        int
	Status      int8
	CreatedAt   int64 `gorm:"autoCreateTime"`
	UpdatedAt   int64 `gorm:"autoUpdateTime"`
}

func (Brand) TableName() string {
	return "brands"
}

type Order struct {
	ID              string `gorm:"primaryKey;size:64"`
	TenantID        int64  `gorm:"index;not null"`
	UserID          int64  `gorm:"index;not null"`
	OrderNo         string `gorm:"size:64;uniqueIndex;not null"`
	Status          int8
	TotalAmount     int64
	DiscountAmount  int64
	FreightAmount   int64
	PayAmount       int64
	Currency        string `gorm:"size:10;default:'CNY'"`
	AddressName     string `gorm:"size:100"`
	AddressPhone    string `gorm:"size:20"`
	AddressProvince string `gorm:"size:50"`
	AddressCity     string `gorm:"size:50"`
	AddressDistrict string `gorm:"size:50"`
	AddressDetail   string `gorm:"type:text"`
	AddressZipcode  string `gorm:"size:20"`
	TrackingNo      string `gorm:"size:100"`
	Carrier         string `gorm:"size:50"`
	Remark          string `gorm:"type:text"`
	ExpireAt        int64
	PaidAt          *int64
	ShippedAt       *int64
	CompletedAt     *int64
	CancelledAt     *int64
	CreatedAt       int64 `gorm:"autoCreateTime"`
	UpdatedAt       int64 `gorm:"autoUpdateTime"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID          int64  `gorm:"primaryKey"`
	OrderID     string `gorm:"index;size:64;not null"`
	ProductID   int64  `gorm:"index;not null"`
	SKUId       int64  `gorm:"not null"`
	ProductName string `gorm:"size:255;not null"`
	SKUName     string `gorm:"size:255"`
	Image       string `gorm:"size:500"`
	Price       int64
	Quantity    int
	TotalAmount int64
	CreatedAt   int64 `gorm:"autoCreateTime"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

type Cart struct {
	ID        int64  `gorm:"primaryKey"`
	TenantID  int64  `gorm:"index;not null"`
	UserID    int64  `gorm:"uniqueIndex:idx_tenant_user"`
	SessionID string `gorm:"size:255;uniqueIndex:idx_tenant_session"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
}

func (Cart) TableName() string {
	return "carts"
}

type CartItem struct {
	ID          int64  `gorm:"primaryKey"`
	TenantID    int64  `gorm:"index;not null"`
	UserID      int64  `gorm:"not null"`
	CartID      int64  `gorm:"index;not null"`
	ProductID   int64  `gorm:"index;not null"`
	SKUId       int64  `gorm:"not null"`
	ProductName string `gorm:"size:255;not null"`
	SKUName     string `gorm:"size:255"`
	Image       string `gorm:"size:500"`
	Price       int64
	Quantity    int
	TotalAmount int64
	Selected    int8
	CreatedAt   int64 `gorm:"autoCreateTime"`
	UpdatedAt   int64 `gorm:"autoUpdateTime"`
}

func (CartItem) TableName() string {
	return "cart_items"
}

type Payment struct {
	ID            string `gorm:"primaryKey;size:64"`
	TenantID      int64  `gorm:"index;not null"`
	OrderID       string `gorm:"index;size:64;not null"`
	UserID        int64  `gorm:"not null"`
	Amount        int64
	Status        int8
	Method        int8
	TransactionID string `gorm:"size:255"`
	PaidAt        *int64
	ExpireAt      int64
	NotifyURL     string `gorm:"size:500"`
	ReturnURL     string `gorm:"size:500"`
	CreatedAt     int64  `gorm:"autoCreateTime"`
	UpdatedAt     int64  `gorm:"autoUpdateTime"`
}

func (Payment) TableName() string {
	return "payments"
}

type Coupon struct {
	ID           int64  `gorm:"primaryKey"`
	TenantID     int64  `gorm:"index;not null"`
	Name         string `gorm:"size:255;not null"`
	Code         string `gorm:"size:100;uniqueIndex;not null"`
	Description  string `gorm:"type:text"`
	Type         int8
	Value        int64
	MinAmount    int64
	MaxDiscount  int64
	TotalCount   int
	UsedCount    int
	PerUserLimit int
	Status       int8
	StartAt      int64
	EndAt        int64
	CreatedAt    int64 `gorm:"autoCreateTime"`
	UpdatedAt    int64 `gorm:"autoUpdateTime"`
}

func (Coupon) TableName() string {
	return "coupons"
}

type UserCoupon struct {
	ID         int64 `gorm:"primaryKey"`
	TenantID   int64 `gorm:"index;not null"`
	UserID     int64 `gorm:"index;not null"`
	CouponID   int64 `gorm:"index;not null"`
	Status     int8
	UsedAt     *int64
	OrderID    string `gorm:"size:64"`
	ReceivedAt int64  `gorm:"autoCreateTime"`
	ExpireAt   int64
}

func (UserCoupon) TableName() string {
	return "user_coupons"
}

type Shop struct {
	ID             int64  `gorm:"primaryKey"`
	TenantID       int64  `gorm:"uniqueIndex;not null"`
	Name           string `gorm:"size:255;not null"`
	Description    string `gorm:"type:text"`
	Logo           string `gorm:"size:500"`
	Banner         string `gorm:"size:500"`
	ContactPhone   string `gorm:"size:20"`
	ContactEmail   string `gorm:"size:255"`
	Address        string `gorm:"type:text"`
	SocialLinks    string `gorm:"type:json"`
	SeoTitle       string `gorm:"size:255"`
	SeoDescription string `gorm:"type:text"`
	SeoKeywords    string `gorm:"size:500"`
	Status         int8
	CreatedAt      int64 `gorm:"autoCreateTime"`
	UpdatedAt      int64 `gorm:"autoUpdateTime"`
}

func (Shop) TableName() string {
	return "shops"
}

type Theme struct {
	ID          int64  `gorm:"primaryKey"`
	TenantID    int64  `gorm:"index;not null"`
	Name        string `gorm:"size:100;not null"`
	Code        string `gorm:"size:100;not null"`
	Description string `gorm:"type:text"`
	Thumbnail   string `gorm:"size:500"`
	Config      string `gorm:"type:json"`
	IsActive    int8
	IsCustom    int8
	CreatedAt   int64 `gorm:"autoCreateTime"`
	UpdatedAt   int64 `gorm:"autoUpdateTime"`
}

func (Theme) TableName() string {
	return "themes"
}
