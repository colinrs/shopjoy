// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	"github.com/colinrs/shopjoy/admin/internal/config"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/pkg/infra"
	"github.com/colinrs/shopjoy/pkg/snowflake"
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
