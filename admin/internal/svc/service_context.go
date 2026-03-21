// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"time"

	appAdminUser "github.com/colinrs/shopjoy/admin/internal/application/adminuser"
	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	"github.com/colinrs/shopjoy/admin/internal/config"
	"github.com/colinrs/shopjoy/admin/internal/domain/market"
	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/middleware"
	"github.com/colinrs/shopjoy/pkg/auth"
	"github.com/colinrs/shopjoy/pkg/infra"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config              config.Config
	DB                  *gorm.DB
	ProductService      appProduct.Service
	UserService         appUser.Service
	AdminUserService    appAdminUser.Service
	JWTManager          *auth.JWTManager
	AuthMiddleware      rest.Middleware
	ProductMarketRepo   product.ProductMarketRepository
	MarketRepo          market.Repository
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

	return &ServiceContext{
		Config:           c,
		DB:               db,
		ProductService:   productService,
		UserService:      userService,
		AdminUserService: adminUserService,
		JWTManager:       jwtManager,
		AuthMiddleware:   authMiddleware,
		ProductMarketRepo: productMarketRepo,
		MarketRepo:        marketRepo,
	}
}
