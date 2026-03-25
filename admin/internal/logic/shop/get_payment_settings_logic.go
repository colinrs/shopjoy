package shop

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaymentSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPaymentSettingsLogic {
	return GetPaymentSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaymentSettingsLogic) GetPaymentSettings() (resp *types.PaymentSettings, err error) {
	// todo: add your logic here and delete this line

	return
}
