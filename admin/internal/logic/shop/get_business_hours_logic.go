package shop

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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

func (l *GetBusinessHoursLogic) GetBusinessHours() (resp []types.BusinessHours, err error) {
	// todo: add your logic here and delete this line

	return
}
