package auth

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/auth"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func RegisterTenantAdminHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterTenantAdminRequest
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := auth.NewRegisterTenantAdminLogic(r.Context(), svcCtx)
		resp, err := l.RegisterTenantAdmin(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
