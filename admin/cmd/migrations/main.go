// Package main provides database migration tool
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/adminuser"
	"github.com/colinrs/shopjoy/admin/internal/domain/cart"
	"github.com/colinrs/shopjoy/admin/internal/domain/coupon"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/market"
	"github.com/colinrs/shopjoy/admin/internal/domain/order"
	"github.com/colinrs/shopjoy/admin/internal/domain/payment"
	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/domain/promotion"
	"github.com/colinrs/shopjoy/admin/internal/domain/review"
	"github.com/colinrs/shopjoy/admin/internal/domain/role"
	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/domain/shop"
	"github.com/colinrs/shopjoy/admin/internal/domain/storefront"
	"github.com/colinrs/shopjoy/admin/internal/domain/tenant"
	"github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dsn         = flag.String("dsn", "root:password@tcp(localhost:3306)/shopjoy?charset=utf8mb4&parseTime=True&loc=UTC", "database DSN")
	migrateOnly = flag.Bool("migrate-only", false, "only run migrations, skip seed data")
	seedOnly    = flag.Bool("seed-only", false, "only run seed data, skip migrations")
	clean       = flag.Bool("clean", false, "drop all tables before migrating (WARNING: destructive)")
)

func main() {
	flag.Parse()

	ctx := context.Background()
	db, err := connect(*dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if *clean {
		log.Println("WARNING: Dropping all tables...")
		if err := dropAll(db); err != nil {
			log.Fatalf("failed to drop tables: %v", err)
		}
	}

	if !*seedOnly {
		if err := migrate(ctx, db); err != nil {
			log.Fatalf("failed to migrate: %v", err)
		}
		log.Println("Migration completed successfully")
	}

	if !*migrateOnly {
		if err := seed(ctx, db); err != nil {
			log.Fatalf("failed to seed: %v", err)
		}
		log.Println("Seed data completed successfully")
	}
}

func connect(dsn string) (*gorm.DB, error) {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying db: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func dropAll(db *gorm.DB) error {
	return db.Exec(`
		DROP TABLE IF EXISTS points_redemptions,
		points_transactions, points_accounts, redeem_rules, earn_rules,
		review_replies, reviews, review_stats,
		shipment_items, shipments, refunds, carriers, refund_reasons,
		shipping_template_mappings, shipping_zone_regions, shipping_zones, shipping_templates,
		payment_transactions, webhook_events, payment_refunds, order_payments,
		order_items, orders,
		coupon_user, coupons, promotions,
		skus, products, categories, brands, markets,
		cart_items, carts,
		admin_users, role_permissions, user_roles, roles, permissions,
		shop_notification_settings, shop_payment_settings, shop_shipping_settings,
		shop_business_hours, shop_settings,
		pages, themes, shops, navigations, nav_items, page_versions,
		decorations, seo_configs, theme_audit_logs,
		tenants, users
	`).Error
}

func migrate(ctx context.Context, db *gorm.DB) error {
	models := []interface{}{
		// Base tables - no foreign key dependencies
		&tenant.Tenant{},
		&user.User{},
		&user.UserAddress{},
		&role.Role{},
		&role.Permission{},
		&role.RolePermission{},
		&role.UserRole{},

		// Admin users
		&adminuser.AdminUser{},

		// Product tables
		&product.Category{},
		&product.Brand{},
		&market.Market{},
		&product.Product{},
		&product.SKU{},

		// Cart
		&cart.Cart{},
		&cart.CartItem{},

		// Promotion tables
		&promotion.Promotion{},
		&promotion.PromotionRule{},
		&coupon.Coupon{},
		&coupon.UserCoupon{},

		// Order tables
		&order.Order{},
		&order.OrderItem{},

		// Payment tables
		&payment.Payment{},
		&payment.PaymentRefund{},
		&payment.PaymentTransaction{},
		&payment.WebhookEvent{},

		// Fulfillment tables
		&fulfillment.Carrier{},
		&fulfillment.RefundReason{},
		&fulfillment.Shipment{},
		&fulfillment.ShipmentItem{},
		&fulfillment.Refund{},

		// Shipping tables
		&shipping.ShippingTemplate{},
		&shipping.ShippingZone{},
		&shipping.ShippingZoneRegion{},
		&shipping.ShippingTemplateMapping{},

		// Storefront tables
		&storefront.Shop{},
		&storefront.Theme{},
		&storefront.Page{},
		&storefront.Decoration{},
		&storefront.PageVersion{},
		&storefront.SEOConfigEntity{},
		&storefront.Navigation{},
		&storefront.NavItem{},
		&storefront.ThemeAuditLog{},

		// Shop settings
		&shop.ShopSettings{},
		&shop.BusinessHours{},
		&shop.NotificationSettings{},
		&shop.PaymentSettings{},
		&shop.ShippingSettings{},

		// Points tables
		&points.EarnRule{},
		&points.RedeemRule{},
		&points.PointsAccount{},
		&points.PointsTransaction{},
		&points.PointsRedemption{},

		// Review tables
		&review.Review{},
		&review.ReviewReply{},
		&review.ReviewStats{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
		log.Printf("Migrated: %T", model)
	}

	return nil
}

func seed(ctx context.Context, db *gorm.DB) error {
	now := time.Now().UTC()

	// Seed tenant
	t := &tenant.Tenant{
		Model:  application.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		Name:   "Demo Tenant",
		Code:   "demo",
		Status: tenant.StatusActive,
		Plan:   tenant.PlanPro,
	}
	if err := db.FirstOrCreate(t, &tenant.Tenant{Code: "demo"}).Error; err != nil {
		return err
	}

	// Seed admin user
	password, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := &adminuser.AdminUser{
		Model:    application.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		TenantID: 0,
		Username: "admin",
		Email:    "admin@shopjoy.com",
		Password: string(password),
		RealName: "Administrator",
		Type:     adminuser.TypeSuperAdmin,
		Status:   adminuser.StatusActive,
	}
	if err := db.FirstOrCreate(admin, &adminuser.AdminUser{Username: "admin"}).Error; err != nil {
		return err
	}

	// Seed tenant admin
	tenantAdmin := &adminuser.AdminUser{
		Model:    application.Model{ID: 2, CreatedAt: now, UpdatedAt: now},
		TenantID: 1,
		Username: "tenant_admin",
		Email:    "tenant@shopjoy.com",
		Password: string(password),
		RealName: "Tenant Admin",
		Type:     adminuser.TypeTenantAdmin,
		Status:   adminuser.StatusActive,
	}
	if err := db.FirstOrCreate(tenantAdmin, &adminuser.AdminUser{Username: "tenant_admin"}).Error; err != nil {
		return err
	}

	// Seed roles
	adminRole := &role.Role{
		Model:       application.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		TenantID:    0,
		Name:        "Super Admin",
		Code:        "super_admin",
		Description: "System super administrator",
		Status:      role.StatusEnabled,
		IsSystem:    true,
	}
	if err := db.FirstOrCreate(adminRole, &role.Role{Code: "super_admin"}).Error; err != nil {
		return err
	}

	tenantRole := &role.Role{
		Model:       application.Model{ID: 2, CreatedAt: now, UpdatedAt: now},
		TenantID:    1,
		Name:        "Tenant Admin",
		Code:        "tenant_admin",
		Description: "Tenant administrator",
		Status:      role.StatusEnabled,
		IsSystem:    false,
	}
	if err := db.FirstOrCreate(tenantRole, &role.Role{Code: "tenant_admin"}).Error; err != nil {
		return err
	}

	// Seed permissions
	permissions := []role.Permission{
		{Model: application.Model{ID: 1, CreatedAt: now, UpdatedAt: now}, Name: "Dashboard", Code: "dashboard", Type: role.PermissionTypeMenu, Path: "/dashboard", Icon: "dashboard", Sort: 1},
		{Model: application.Model{ID: 2, CreatedAt: now, UpdatedAt: now}, Name: "Products", Code: "products", Type: role.PermissionTypeMenu, Path: "/products", Icon: "product", Sort: 2},
		{Model: application.Model{ID: 3, CreatedAt: now, UpdatedAt: now}, Name: "Orders", Code: "orders", Type: role.PermissionTypeMenu, Path: "/orders", Icon: "order", Sort: 3},
		{Model: application.Model{ID: 4, CreatedAt: now, UpdatedAt: now}, Name: "Customers", Code: "customers", Type: role.PermissionTypeMenu, Path: "/customers", Icon: "user", Sort: 4},
		{Model: application.Model{ID: 5, CreatedAt: now, UpdatedAt: now}, Name: "Promotions", Code: "promotions", Type: role.PermissionTypeMenu, Path: "/promotions", Icon: "promotion", Sort: 5},
		{Model: application.Model{ID: 6, CreatedAt: now, UpdatedAt: now}, Name: "Settings", Code: "settings", Type: role.PermissionTypeMenu, Path: "/settings", Icon: "setting", Sort: 6},
	}
	for _, p := range permissions {
		db.FirstOrCreate(&p, &role.Permission{Code: p.Code})
	}

	// Assign all permissions to super_admin role
	for _, p := range permissions {
		rp := &role.RolePermission{
			Model:        application.Model{CreatedAt: now, UpdatedAt: now},
			RoleID:       adminRole.ID,
			PermissionID: p.ID,
		}
		db.FirstOrCreate(rp, &role.RolePermission{RoleID: adminRole.ID, PermissionID: p.ID})
	}

	// Assign tenant_admin role to admin user
	ur1 := &role.UserRole{
		Model:  application.Model{CreatedAt: now, UpdatedAt: now},
		UserID: admin.ID,
		RoleID: adminRole.ID,
	}
	db.FirstOrCreate(ur1, &role.UserRole{UserID: admin.ID, RoleID: adminRole.ID})

	// Assign tenant_admin role to tenant_admin user
	ur2 := &role.UserRole{
		Model:  application.Model{CreatedAt: now, UpdatedAt: now},
		UserID: tenantAdmin.ID,
		RoleID: tenantRole.ID,
	}
	db.FirstOrCreate(ur2, &role.UserRole{UserID: tenantAdmin.ID, RoleID: tenantRole.ID})

	// Seed markets
	markets := []market.Market{
		{Model: application.Model{ID: 1, CreatedAt: now, UpdatedAt: now}, Code: "US", Name: "United States", Currency: "USD", DefaultLanguage: "en", Flag: "🇺🇸", IsActive: true, IsDefault: true},
		{Model: application.Model{ID: 2, CreatedAt: now, UpdatedAt: now}, Code: "CN", Name: "China", Currency: "CNY", DefaultLanguage: "zh", Flag: "🇨🇳", IsActive: true, IsDefault: false},
		{Model: application.Model{ID: 3, CreatedAt: now, UpdatedAt: now}, Code: "UK", Name: "United Kingdom", Currency: "GBP", DefaultLanguage: "en", Flag: "🇬🇧", IsActive: true, IsDefault: false},
	}
	for _, m := range markets {
		db.FirstOrCreate(&m, &market.Market{Code: m.Code})
	}

	// Seed categories
	categories := []product.Category{
		{Model: application.Model{ID: 1, CreatedAt: now, UpdatedAt: now}, TenantID: 1, ParentID: 0, Name: "Electronics", Code: "electronics", Sort: 1, Status: product.CategoryStatusEnabled},
		{Model: application.Model{ID: 2, CreatedAt: now, UpdatedAt: now}, TenantID: 1, ParentID: 0, Name: "Clothing", Code: "clothing", Sort: 2, Status: product.CategoryStatusEnabled},
		{Model: application.Model{ID: 3, CreatedAt: now, UpdatedAt: now}, TenantID: 1, ParentID: 1, Name: "Smartphones", Code: "smartphones", Sort: 1, Status: product.CategoryStatusEnabled},
		{Model: application.Model{ID: 4, CreatedAt: now, UpdatedAt: now}, TenantID: 1, ParentID: 1, Name: "Laptops", Code: "laptops", Sort: 2, Status: product.CategoryStatusEnabled},
	}
	for _, c := range categories {
		db.FirstOrCreate(&c, &product.Category{Code: c.Code, TenantID: c.TenantID})
	}

	// Seed brands
	brands := []product.Brand{
		{Model: application.Model{ID: 1, CreatedAt: now, UpdatedAt: now}, TenantID: 1, Name: "Apple", Status: shared.StatusEnabled},
		{Model: application.Model{ID: 2, CreatedAt: now, UpdatedAt: now}, TenantID: 1, Name: "Samsung", Status: shared.StatusEnabled},
		{Model: application.Model{ID: 3, CreatedAt: now, UpdatedAt: now}, TenantID: 1, Name: "Nike", Status: shared.StatusEnabled},
	}
	for _, b := range brands {
		db.FirstOrCreate(&b, &product.Brand{Name: b.Name, TenantID: b.TenantID})
	}

	// Seed products
	products := []*product.Product{
		{
			Model:           application.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
			TenantID:        1,
			SKU:             "IPHONE15-256",
			Name:            "iPhone 15 Pro 256GB",
			Description:     "Latest Apple smartphone with A17 chip",
			Price:           product.Money{Amount: decimal.NewFromInt(999), Currency: "USD"},
			CostPrice:       product.Money{Amount: decimal.NewFromInt(700), Currency: "USD"},
			Stock:           100,
			Status:          product.StatusOnSale,
			CategoryID:      3,
			Brand:           "Apple",
			SKUPrefix:       "IPHONE",
			IsMatrixProduct: false,
		},
		{
			Model:           application.Model{ID: 2, CreatedAt: now, UpdatedAt: now},
			TenantID:        1,
			SKU:             "GALAXY-S24-512",
			Name:            "Samsung Galaxy S24 Ultra",
			Description:     "Premium Android smartphone",
			Price:           product.Money{Amount: decimal.NewFromInt(899), Currency: "USD"},
			CostPrice:       product.Money{Amount: decimal.NewFromInt(650), Currency: "USD"},
			Stock:           80,
			Status:          product.StatusOnSale,
			CategoryID:      3,
			Brand:           "Samsung",
			SKUPrefix:       "GALAXY",
			IsMatrixProduct: false,
		},
		{
			Model:           application.Model{ID: 3, CreatedAt: now, UpdatedAt: now},
			TenantID:        1,
			SKU:             "NIKE-AIR-MAX",
			Name:            "Nike Air Max 90",
			Description:     "Classic running shoes",
			Price:           product.Money{Amount: decimal.NewFromInt(129), Currency: "USD"},
			CostPrice:       product.Money{Amount: decimal.NewFromInt(60), Currency: "USD"},
			Stock:           200,
			Status:          product.StatusOnSale,
			CategoryID:      2,
			Brand:           "Nike",
			SKUPrefix:       "NIKE",
			IsMatrixProduct: false,
		},
	}
	for _, p := range products {
		db.FirstOrCreate(p, &product.Product{TenantID: p.TenantID, SKU: p.SKU})
	}

	// Seed storefront shop
	shop := &storefront.Shop{
		Model:       application.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		TenantID:    1,
		Name:        "ShopJoy Store",
		Description: "Your one-stop shop for everything",
		Status:      shared.StatusEnabled,
	}
	if err := db.FirstOrCreate(shop, &storefront.Shop{TenantID: 1}).Error; err != nil {
		return err
	}

	// Seed theme
	theme := &storefront.Theme{
		Model:       application.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		TenantID:    1,
		Name:        "Default Theme",
		Code:        "default",
		Description: "Clean and modern default theme",
		IsActive:    true,
		IsPreset:    true,
	}
	if err := db.FirstOrCreate(theme, &storefront.Theme{Code: "default"}).Error; err != nil {
		return err
	}

	// Seed carriers
	for _, c := range fulfillment.DefaultCarriers {
		c.Model = application.Model{ID: 0, CreatedAt: now, UpdatedAt: now}
		db.FirstOrCreate(&c, &fulfillment.Carrier{Code: c.Code})
	}

	// Seed refund reasons
	for _, r := range fulfillment.DefaultRefundReasons {
		r.Model = application.Model{ID: 0, CreatedAt: now, UpdatedAt: now}
		db.FirstOrCreate(&r, &fulfillment.RefundReason{Code: r.Code})
	}

	// Seed points earn rules
	earnRule := &points.EarnRule{
		Model:            application.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		TenantID:         1,
		Name:             "Order Payment Points",
		Description:      "Earn 1 point per dollar spent",
		Scenario:         points.EarnScenarioOrderPayment,
		CalculationType:  points.CalculationTypeRatio,
		Ratio:            decimal.NewFromInt(1),
		ConditionType:    points.ConditionTypeNone,
		ExpirationMonths: 12,
		Status:           points.EarnRuleStatusActive,
	}
	if err := db.FirstOrCreate(earnRule, &points.EarnRule{TenantID: 1, Scenario: points.EarnScenarioOrderPayment}).Error; err != nil {
		return err
	}

	signInRule := &points.EarnRule{
		Model:            application.Model{ID: 2, CreatedAt: now, UpdatedAt: now},
		TenantID:         1,
		Name:             "Daily Sign In",
		Description:      "Earn 10 points for daily sign in",
		Scenario:         points.EarnScenarioSignIn,
		CalculationType:  points.CalculationTypeFixed,
		FixedPoints:      10,
		ConditionType:    points.ConditionTypeNone,
		ExpirationMonths: 30,
		Status:           points.EarnRuleStatusActive,
	}
	if err := db.FirstOrCreate(signInRule, &points.EarnRule{TenantID: 1, Scenario: points.EarnScenarioSignIn}).Error; err != nil {
		return err
	}

	log.Println("Seed data completed successfully")
	return nil
}
