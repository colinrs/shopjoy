package users

import (
	"context"
	"net/http"
	"time"

	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	domain "github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
	w      http.ResponseWriter
}

func NewExportUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) ExportUsersLogic {
	return ExportUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
		r:      r,
	}
}

func (l *ExportUsersLogic) ExportUsers(req *types.ExportUsersRequest) error {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		return code.ErrTenantInvalidID
	}

	queryReq := appUser.EnhancedQueryRequest{
		PageQuery: shared.PageQuery{
			Page:     1,
			PageSize: 10001, // Fetch one more to check if limit exceeded
		},
		Keyword:       req.Keyword,
		Status:        domain.Status(req.Status),
		RegisterStart: req.RegisterStart,
		RegisterEnd:   req.RegisterEnd,
	}

	// Get users for export
	data, err := l.svcCtx.UserService.ExportUsers(l.ctx, tenantID, queryReq)
	if err != nil {
		return err
	}

	// Set response headers for CSV download
	l.w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	l.w.Header().Set("Content-Disposition", "attachment; filename=users_"+time.Now().Format("20060102150405")+".csv")

	// Write UTF-8 BOM for Excel compatibility
	l.w.Write([]byte{0xEF, 0xBB, 0xBF})

	// Write the CSV data returned from service
	_, err = l.w.Write(data)
	return err
}
