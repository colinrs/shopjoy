// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"time"

	appAdminUser "github.com/colinrs/shopjoy/admin/internal/application/adminuser"
	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	"github.com/colinrs/shopjoy/admin/internal/config"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/market"
	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
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
	userService := appUser.NewService(db, userRepo, idGen)

	adminUserRepo := persistence.NewAdminUserRepository()
	adminUserService := appAdminUser.NewService(adminUserRepo, db, c.JWT.Secret)

	jwtManager := auth.NewJWTManager(c.JWT.Secret, time.Duration(c.JWT.AccessExpiry)*time.Second, time.Duration(c.JWT.RefreshExpiry)*time.Second)

	authMiddleware := middleware.NewAuthMiddleware(c.JWT.Secret)

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
	// Use NopOrderValidator until order service integration is complete
	orderFulfillmentApp := appfulfillment.NewOrderFulfillmentApp(db, shipmentRepo, shipmentItemRepo, carrierRepo, refundRepo, orderRepo, orderItemRepo, idGen, &appfulfillment.NopOrderValidator{})
	refundApp := appfulfillment.NewRefundApp(db, refundRepo, refundReasonRepo, idGen)
	refundReasonApp := appfulfillment.NewRefundReasonApp(db, refundReasonRepo)

	return &ServiceContext{
		Config:                 c,
		DB:                     db,
		ProductService:         productService,
		UserService:            userService,
		AdminUserService:       adminUserService,
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
	}
}