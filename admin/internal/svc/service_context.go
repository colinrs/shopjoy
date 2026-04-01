// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"time"

	appAdminUser "github.com/colinrs/shopjoy/admin/internal/application/adminuser"
	"github.com/colinrs/shopjoy/admin/internal/domain/adminuser"
	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	appPayment "github.com/colinrs/shopjoy/admin/internal/application/payment"
	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	appPoints "github.com/colinrs/shopjoy/admin/internal/application/points"
	appReview "github.com/colinrs/shopjoy/admin/internal/application/review"
	appStorefront "github.com/colinrs/shopjoy/admin/internal/application/storefront"
	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	"github.com/colinrs/shopjoy/admin/internal/config"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/market"
	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/domain/review"
	"github.com/colinrs/shopjoy/admin/internal/domain/role"
	"github.com/colinrs/shopjoy/admin/internal/domain/shop"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/storage"
	"github.com/colinrs/shopjoy/admin/internal/middleware"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/auth"
	"github.com/colinrs/shopjoy/pkg/infra"
	"github.com/colinrs/shopjoy/pkg/sku"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config                 config.Config
	DB                     *gorm.DB
	ProductService         appProduct.Service
	UserService            appUser.Service
	AdminUserService       appAdminUser.Service
	PaymentService         appPayment.Service
	PromotionApp           apppromotion.PromotionApp
	CouponApp              apppromotion.CouponApp
	ShipmentApp            appfulfillment.ShipmentApp
	CarrierApp             appfulfillment.CarrierApp
	OrderFulfillmentApp    appfulfillment.OrderFulfillmentApp
	RefundApp              appfulfillment.RefundApp
	RefundReasonApp        appfulfillment.RefundReasonApp
	JWTManager             *auth.JWTManager
	AuthMiddleware         rest.Middleware
	ProductRepo            product.Repository
	ProductMarketRepo      product.ProductMarketRepository
	MarketRepo             market.Repository
	CategoryRepo           product.CategoryRepository
	CategoryMarketRepo     product.CategoryMarketRepository
	BrandRepo              product.BrandRepository
	BrandMarketRepo        product.BrandMarketRepository
	WarehouseRepo          product.WarehouseRepository
	WarehouseInventoryRepo product.WarehouseInventoryRepository
	InventoryLogRepo       product.InventoryLogRepository
	SKURepo                   product.SKURepository
	ProductLocalizationRepo   product.ProductLocalizationRepository
	SKUGenerator              sku.Generator
	IDGen                     snowflake.Snowflake
	PromotionRepo             pkgpromotion.Repository
	CouponRepo                pkgpromotion.CouponRepository
	UserCouponRepo            pkgpromotion.UserCouponRepository
	OrderRepo                 fulfillment.OrderRepository
	OrderItemRepo             fulfillment.OrderItemRepository
	ReviewService             appReview.Service
	ReviewRepo                review.Repository
	ReplyRepo                 review.ReplyRepository
	// Storefront services
	ThemeService      appStorefront.ThemeService
	PageService       appStorefront.PageService
	DecorationService appStorefront.DecorationService
	VersionService    appStorefront.VersionService
	SEOService        appStorefront.SEOService
	// Role and Permission
	RoleRepo       role.Repository
	PermissionRepo role.PermissionRepository
	AdminUserRepo  adminuser.Repository
	// Shipping
	ShippingRepo persistence.ShippingTemplateRepository
	// Storage
	Storage storage.Storage
	// Points
	PointsService          appPoints.Service
	EarnRuleRepo           points.EarnRuleRepository
	RedeemRuleRepo         points.RedeemRuleRepository
	PointsAccountRepo      points.PointsAccountRepository
	PointsTransactionRepo  points.PointsTransactionRepository
	PointsRedemptionRepo   points.PointsRedemptionRepository
	// Shop Settings
	ShopSettingsRepo         shop.ShopSettingsRepository
	BusinessHoursRepo        shop.BusinessHoursRepository
	NotificationSettingsRepo shop.NotificationSettingsRepository
	PaymentSettingsRepo     shop.PaymentSettingsRepository
	ShippingSettingsRepo     shop.ShippingSettingsRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := infra.Database(&c.MySQL)
	if err != nil {
		logx.Errorf("init database error: %v", err)
		panic(err)
	}

	idGen := snowflake.NewSnowflake(1)

	productRepo := persistence.NewProductRepository()
	productService := appProduct.NewService(db, productRepo, idGen)

	userRepo := persistence.NewUserRepository()
	addressRepo := persistence.NewUserAddressRepository()
	userService := appUser.NewService(db, userRepo, addressRepo, idGen)

	adminUserRepo := persistence.NewAdminUserRepository()
	roleRepo := persistence.NewRoleRepository()
	permissionRepo := persistence.NewPermissionRepository()
	adminUserService := appAdminUser.NewService(adminUserRepo, roleRepo, db, c.JWT.Secret)

	jwtManager := auth.NewJWTManager(c.JWT.Secret, time.Duration(c.JWT.AccessExpiry)*time.Second, time.Duration(c.JWT.RefreshExpiry)*time.Second)

	authMiddleware := middleware.NewAuthMiddleware(c.JWT.Secret, db, adminUserRepo)

	productMarketRepo := persistence.NewProductMarketRepository()
	marketRepo := persistence.NewMarketRepository()
	categoryRepo := persistence.NewCategoryRepository()
	categoryMarketRepo := persistence.NewCategoryMarketRepository()
	brandRepo := persistence.NewBrandRepository()
	brandMarketRepo := persistence.NewBrandMarketRepository()
	warehouseRepo := persistence.NewWarehouseRepository()
	warehouseInventoryRepo := persistence.NewWarehouseInventoryRepository()
	inventoryLogRepo := persistence.NewInventoryLogRepository()
	skuRepo := persistence.NewSKURepository()
	productLocalizationRepo := persistence.NewProductLocalizationRepository()
	skuGenerator := sku.NewGenerator(sku.DefaultConfig())

	// Promotion and Coupon repositories
	promotionRepo := persistence.NewPromotionRepository()
	couponRepo := persistence.NewCouponRepository()
	userCouponRepo := persistence.NewUserCouponRepository()

	// Promotion and Coupon application services
	promotionApp := apppromotion.NewPromotionApp(db, promotionRepo, idGen)
	couponApp := apppromotion.NewCouponApp(db, couponRepo, userCouponRepo, idGen)

	// Fulfillment repositories
	shipmentRepo := persistence.NewShipmentRepository()
	shipmentItemRepo := persistence.NewShipmentItemRepository()
	carrierRepo := persistence.NewCarrierRepository()
	refundRepo := persistence.NewRefundRepository()
	refundReasonRepo := persistence.NewRefundReasonRepository()
	orderRepo := persistence.NewOrderRepository()
	orderItemRepo := persistence.NewOrderItemRepository()

	// Fulfillment application services
	shipmentApp := appfulfillment.NewShipmentApp(db, shipmentRepo, shipmentItemRepo, carrierRepo, idGen)
	carrierApp := appfulfillment.NewCarrierApp(db, carrierRepo)
	// Use DefaultOrderValidator until order service integration is complete
	orderFulfillmentApp := appfulfillment.NewOrderFulfillmentApp(db, shipmentRepo, shipmentItemRepo, carrierRepo, refundRepo, orderRepo, orderItemRepo, idGen, &appfulfillment.DefaultOrderValidator{})
	refundApp := appfulfillment.NewRefundApp(db, refundRepo, refundReasonRepo, idGen)
	refundReasonApp := appfulfillment.NewRefundReasonApp(db, refundReasonRepo)

	// Review repositories
	reviewRepo := persistence.NewReviewRepository()
	replyRepo := persistence.NewReplyRepository()

	// Review service
	reviewService := appReview.NewService(db, reviewRepo, replyRepo, idGen)

	// Payment repositories
	paymentRepo := persistence.NewPaymentRepository()
	paymentRefundRepo := persistence.NewPaymentRefundRepository()
	paymentTransactionRepo := persistence.NewPaymentTransactionRepository()
	webhookEventRepo := persistence.NewWebhookEventRepository()
	paymentSettingsRepo := persistence.NewPaymentSettingsRepository()

	// Payment service
	paymentService := appPayment.NewService(db, paymentRepo, paymentRefundRepo, paymentTransactionRepo, webhookEventRepo, paymentSettingsRepo, idGen)

	// Storefront repositories
	themeRepo := persistence.NewThemeRepository()
	pageRepo := persistence.NewPageRepository()
	decorationRepo := persistence.NewDecorationRepository()
	pageVersionRepo := persistence.NewPageVersionRepository()
	seoConfigRepo := persistence.NewSEOConfigRepository()
	shopRepo := persistence.NewShopRepository()
	themeAuditLogRepo := persistence.NewThemeAuditLogRepository()

	// Storefront services
	themeService := appStorefront.NewThemeService(db, themeRepo, shopRepo, themeAuditLogRepo)
	pageService := appStorefront.NewPageService(db, pageRepo, decorationRepo, pageVersionRepo, idGen)
	decorationService := appStorefront.NewDecorationService(db, decorationRepo, idGen)
	versionService := appStorefront.NewVersionService(db, pageRepo, pageVersionRepo, decorationRepo, idGen)
	seoService := appStorefront.NewSEOService(db, seoConfigRepo)

	// Points repositories
	earnRuleRepo := persistence.NewEarnRuleRepository()
	redeemRuleRepo := persistence.NewRedeemRuleRepository()
	pointsAccountRepo := persistence.NewPointsAccountRepository()
	pointsTransactionRepo := persistence.NewPointsTransactionRepository()
	pointsRedemptionRepo := persistence.NewPointsRedemptionRepository()

	// Points service
	pointsService := appPoints.NewService(
		db,
		earnRuleRepo,
		redeemRuleRepo,
		pointsAccountRepo,
		pointsTransactionRepo,
		pointsRedemptionRepo,
		idGen,
	)

	// Shop settings repositories
	shopSettingsRepo := persistence.NewShopSettingsRepository()
	businessHoursRepo := persistence.NewBusinessHoursRepository()
	notificationSettingsRepo := persistence.NewNotificationSettingsRepository()
	shippingSettingsRepo := persistence.NewShippingSettingsRepository()

	return &ServiceContext{
		Config:                 c,
		DB:                     db,
		ProductService:         productService,
		UserService:            userService,
		AdminUserService:       adminUserService,
		PaymentService:         paymentService,
		PromotionApp:           promotionApp,
		CouponApp:              couponApp,
		ShipmentApp:            shipmentApp,
		CarrierApp:             carrierApp,
		OrderFulfillmentApp:    orderFulfillmentApp,
		RefundApp:              refundApp,
		RefundReasonApp:        refundReasonApp,
		JWTManager:             jwtManager,
		AuthMiddleware:         authMiddleware,
		ProductRepo:            productRepo,
		ProductMarketRepo:      productMarketRepo,
		MarketRepo:             marketRepo,
		CategoryRepo:           categoryRepo,
		CategoryMarketRepo:     categoryMarketRepo,
		BrandRepo:              brandRepo,
		BrandMarketRepo:        brandMarketRepo,
		WarehouseRepo:          warehouseRepo,
		WarehouseInventoryRepo: warehouseInventoryRepo,
		InventoryLogRepo:       inventoryLogRepo,
		SKURepo:                skuRepo,
		ProductLocalizationRepo: productLocalizationRepo,
		SKUGenerator:           skuGenerator,
		IDGen:                  idGen,
		PromotionRepo:          promotionRepo,
		CouponRepo:             couponRepo,
		UserCouponRepo:         userCouponRepo,
		OrderRepo:              orderRepo,
		OrderItemRepo:          orderItemRepo,
		ReviewService:          reviewService,
		ReviewRepo:             reviewRepo,
		ReplyRepo:              replyRepo,
		// Storefront services
		ThemeService:      themeService,
		PageService:       pageService,
		DecorationService: decorationService,
		VersionService:    versionService,
		SEOService:        seoService,
		// Role and Permission
		RoleRepo:       roleRepo,
		PermissionRepo: permissionRepo,
		AdminUserRepo:  adminUserRepo,
		// Shipping
		ShippingRepo: persistence.NewShippingTemplateRepository(),

		// Storage
		Storage: storage.MustNewStorage(storage.Config{
			Type: storage.StorageTypeLocal,
		}),

		// Points repositories and service
		PointsService:         pointsService,
		EarnRuleRepo:          earnRuleRepo,
		RedeemRuleRepo:        redeemRuleRepo,
		PointsAccountRepo:     pointsAccountRepo,
		PointsTransactionRepo: pointsTransactionRepo,
		PointsRedemptionRepo:  pointsRedemptionRepo,
		// Shop settings
		ShopSettingsRepo:         shopSettingsRepo,
		BusinessHoursRepo:        businessHoursRepo,
		NotificationSettingsRepo: notificationSettingsRepo,
		PaymentSettingsRepo:     paymentSettingsRepo,
		ShippingSettingsRepo:     shippingSettingsRepo,
	}
}