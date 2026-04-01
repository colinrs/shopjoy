package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteReplyLogic {
	return DeleteReplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteReplyLogic) DeleteReply(req *types.DeleteReplyReq) (resp *types.DeleteReplyResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	if err := l.svcCtx.ReviewService.DeleteReply(l.ctx, shared.TenantID(tenantID), req.ID); err != nil {
		return nil, err
	}

	return &types.DeleteReplyResp{
		Success: true,
	}, nil
}