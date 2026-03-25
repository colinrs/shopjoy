package shop

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePaymentSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePaymentSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdatePaymentSettingsLogic {
	return UpdatePaymentSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePaymentSettingsLogic) UpdatePaymentSettings(req *types.UpdatePaymentSettingsRequest) (resp *types.PaymentSettings, err error) {
	// todo: add your logic here and delete this line

	return
}
