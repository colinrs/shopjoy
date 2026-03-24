package users

import (
	"context"

	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type SuspendUserWithReasonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSuspendUserWithReasonLogic(ctx context.Context, svcCtx *svc.ServiceContext) SuspendUserWithReasonLogic {
	return SuspendUserWithReasonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SuspendUserWithReasonLogic) SuspendUserWithReason(req *types.SuspendUserWithReasonRequest) (resp *types.GetUserResponse, err error) {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		tenantID = shared.TenantID(1) // 默认租户
	}

	suspendReq := appUser.SuspendUserRequest{
		TenantID: tenantID,
		UserID:   req.ID,
		Reason:   req.Reason,
	}

	if err := l.svcCtx.UserService.SuspendWithReason(l.ctx, suspendReq); err != nil {
		return nil, err
	}

	userResp, err := l.svcCtx.UserService.GetByID(l.ctx, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	return toGetUserResponse(userResp), nil
}
