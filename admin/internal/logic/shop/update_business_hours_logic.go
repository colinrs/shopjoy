package shop

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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
	// todo: add your logic here and delete this line

	return nil
}
