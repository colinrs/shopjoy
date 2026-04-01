package users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) BatchUpdateUserStatusLogic {
	return BatchUpdateUserStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchUpdateUserStatusLogic) BatchUpdateUserStatus(req *types.BatchUpdateUserStatusReq) (resp *types.BatchUpdateUserStatusResp, err error) {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		return nil, code.ErrTenantInvalidID
	}

	resp = &types.BatchUpdateUserStatusResp{
		Success: make([]int64, 0, len(req.UserIDs)),
		Failed:  make([]types.BatchStatusFail, 0),
	}

	for _, userID := range req.UserIDs {
		var updateErr error
		switch req.Status {
		case 1: // activate
			updateErr = l.svcCtx.UserService.Activate(l.ctx, tenantID, userID)
		case 2: // suspend
			updateErr = l.svcCtx.UserService.Suspend(l.ctx, tenantID, userID)
		default:
			resp.Failed = append(resp.Failed, types.BatchStatusFail{
				UserID:  userID,
				Code:    40000, // Bad Request
				Message: "invalid status, must be 1 (activate) or 2 (suspend)",
			})
			continue
		}

		if updateErr != nil {
			resp.Failed = append(resp.Failed, types.BatchStatusFail{
				UserID:  userID,
				Code:    getUserErrorCode(updateErr),
				Message: updateErr.Error(),
			})
			continue
		}

		resp.Success = append(resp.Success, userID)
	}

	return resp, nil
}

// getUserErrorCode extracts the error code from an error
func getUserErrorCode(err error) int {
	if err == nil {
		return 0
	}

	// Check for known user errors using errors.Is
	switch {
	case err == code.ErrUserNotFound:
		return code.ErrUserNotFound.Code
	case err == code.ErrUserAlreadyDeleted:
		return code.ErrUserAlreadyDeleted.Code
	case err == code.ErrUserSuspended:
		return code.ErrUserSuspended.Code
	default:
		return code.ErrInternalServer.Code
	}
}
