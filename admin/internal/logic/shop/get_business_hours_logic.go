package shop

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/shop"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetBusinessHoursLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBusinessHoursLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetBusinessHoursLogic {
	return GetBusinessHoursLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBusinessHoursLogic) GetBusinessHours() (resp []*types.BusinessHours, err error) {
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok || tenantID == 0 {
		return nil, code.ErrTenantInvalidID
	}

	// Find shop settings first to get shop ID
	settings, err := l.svcCtx.ShopSettingsRepo.FindByTenantID(l.ctx, l.svcCtx.DB, tenantID)
	if err != nil {
		l.Logger.Errorf("find shop settings error: %v", err)
		return nil, code.ErrInternalServer
	}

	if settings == nil {
		return nil, code.ErrShopNotFound
	}

	hours, err := l.svcCtx.BusinessHoursRepo.FindByShopID(l.ctx, l.svcCtx.DB, settings.ID)
	if err != nil {
		l.Logger.Errorf("get business hours error: %v", err)
		return nil, code.ErrInternalServer
	}

	// Return empty array if no hours set
	if len(hours) == 0 {
		return []*types.BusinessHours{}, nil
	}

	response := make([]*types.BusinessHours, len(hours))
	for i, h := range hours {
		response[i] = toBusinessHoursResponse(h)
	}

	return response, nil
}

func toBusinessHoursResponse(h *shop.BusinessHours) *types.BusinessHours {
	return &types.BusinessHours{
		DayOfWeek: h.DayOfWeek,
		OpenTime:  h.OpenTime,
		CloseTime: h.CloseTime,
		IsClosed:  h.IsClosed,
	}
}
