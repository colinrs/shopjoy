// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/colinrs/shopjoy/pkg/infra"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	appProduct "github.com/colinrs/shopjoy/shop/internal/application/product"
	"github.com/colinrs/shopjoy/shop/internal/config"
	"github.com/colinrs/shopjoy/shop/internal/infrastructure/persistence"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config         config.Config
	DB             *gorm.DB
	ProductService appProduct.Service
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

	return &ServiceContext{
		Config:         c,
		DB:             db,
		ProductService: productService,
	}
}
