package users

import (
	"context"

	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

// recordOperationLog writes one user_operation_logs row via the application
// service. The service's Record method never returns an error and never blocks
// the parent business operation — call this unconditionally after any
// successful user state change.
//
// IP and User-Agent are not currently available from context (no helper
// exists in pkg/contextx); left empty until request metadata propagation is
// added to the auth middleware. The DB columns tolerate empty strings.
//
// Centralized here so all 7 埋点 points (Tasks 9-15) share one source of
// truth for operator identity extraction.
func recordOperationLog(
	ctx context.Context,
	svcCtx *svc.ServiceContext,
	tenantID shared.TenantID,
	userID int64,
	action string,
	reason string,
) {
	svcCtx.OperationLogService.Record(ctx, svcCtx.DB, appUser.RecordOperationLogInput{
		TenantID:     tenantID,
		UserID:       userID,
		Action:       action,
		OperatorID:   contextx.GetCurrentUserID(ctx),
		OperatorName: contextx.GetCurrentUserName(ctx),
		Reason:       reason,
		// IPAddress:  "", // TODO: wire when contextx exposes GetClientIP
		// UserAgent:  "", // TODO: wire when contextx exposes GetUserAgent
	})
}