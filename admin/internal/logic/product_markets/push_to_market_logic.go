// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product_markets

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

type PushToMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 推送商品到市场
func NewPushToMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushToMarketLogic {
	return &PushToMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PushToMarketLogic) PushToMarket(req *types.PushToMarketReq) (resp *types.PushToMarketResp, err error) {
	db := l.svcCtx.DB
	repo := persistence.NewProductMarketRepository()
	marketRepo := persistence.NewMarketRepository()

	// Get tenant ID from context
	tenantID := l.ctx.Value("tenant_id")
	var tid int64
	if tidVal, ok := tenantID.(int64); ok {
		tid = tidVal
	}

	var success, failed []int64

	for i, marketID := range req.MarketIDs {
		// Validate market exists
		market, err := marketRepo.FindByID(l.ctx, db, marketID)
		if err != nil || !market.IsActive {
			failed = append(failed, marketID)
			continue
		}

		// Check if already exists
		existing, _ := repo.FindByProductAndMarket(l.ctx, db, req.ProductID, marketID, nil)
		if existing != nil {
			failed = append(failed, marketID)
			continue
		}

		// Parse price
		var price decimal.Decimal
		if i < len(req.Prices) {
			price, _ = decimal.NewFromString(req.Prices[i])
		}

		// Create ProductMarket
		pm := &product.ProductMarket{
			TenantID:  tid,
			ProductID: req.ProductID,
			MarketID:  marketID,
			Price:     price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := repo.Create(l.ctx, db, pm); err != nil {
			failed = append(failed, marketID)
			continue
		}

		success = append(success, marketID)
	}

	return &types.PushToMarketResp{
		Success: success,
		Failed:  failed,
	}, nil
}
