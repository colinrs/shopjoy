package fulfillment

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"gorm.io/gorm"
)

// CarrierApp 物流公司应用服务接口
type CarrierApp interface {
	ListCarriers(ctx context.Context) ([]*CarrierResponse, error)
}

type carrierApp struct {
	db          *gorm.DB
	carrierRepo fulfillment.CarrierRepository
}

// NewCarrierApp 创建物流公司应用服务
func NewCarrierApp(db *gorm.DB, carrierRepo fulfillment.CarrierRepository) CarrierApp {
	return &carrierApp{
		db:          db,
		carrierRepo: carrierRepo,
	}
}

func (a *carrierApp) ListCarriers(ctx context.Context) ([]*CarrierResponse, error) {
	carriers, err := a.carrierRepo.FindActive(ctx, a.db)
	if err != nil {
		return nil, err
	}

	list := make([]*CarrierResponse, len(carriers))
	for i, c := range carriers {
		list[i] = &CarrierResponse{
			ID:          c.ID,
			Code:        c.Code,
			Name:        c.Name,
			TrackingURL: c.TrackingURL,
			IsActive:    c.IsActive,
			Sort:        c.Sort,
		}
	}

	return list, nil
}
