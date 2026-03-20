package admin_users

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/admin_users"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func DisableAdminUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminUserIDRequest
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := admin_users.NewDisableAdminUserLogic(r.Context(), svcCtx)
		err := l.DisableAdminUser(&req)
		httpy.ResultCtx(r, w, nil, err)
	}
}
