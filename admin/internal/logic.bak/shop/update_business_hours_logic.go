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

type UpdateBusinessHoursLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBusinessHoursLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateBusinessHoursLogic {
	return UpdateBusinessHoursLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBusinessHoursLogic) UpdateBusinessHours(req *types.UpdateBusinessHoursRequest) error {
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok || tenantID == 0 {
		return code.ErrTenantInvalidID
	}

	// Find shop settings first to get shop ID
	settings, err := l.svcCtx.ShopSettingsRepo.FindByTenantID(l.ctx, l.svcCtx.DB, tenantID)
	if err != nil {
		l.Logger.Errorf("find shop settings error: %v", err)
		return code.ErrInternalServer
	}

	if settings == nil {
		return code.ErrShopNotFound
	}

	// Convert request to domain entities
	hours := make([]*shop.BusinessHours, len(req.Hours))
	for i, h := range req.Hours {
		hours[i] = &shop.BusinessHours{
			DayOfWeek: h.DayOfWeek,
			OpenTime:  h.OpenTime,
			CloseTime: h.CloseTime,
			IsClosed:  h.IsClosed,
		}
	}

	if err := l.svcCtx.BusinessHoursRepo.SaveBatch(l.ctx, l.svcCtx.DB, settings.ID, hours); err != nil {
		l.Logger.Errorf("save business hours error: %v", err)
		return code.ErrInternalServer
	}

	return nil
}
