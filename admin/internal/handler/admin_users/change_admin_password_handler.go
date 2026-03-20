package admin_users

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/admin_users"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func ChangeAdminPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminChangePasswordRequest
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := admin_users.NewChangeAdminPasswordLogic(r.Context(), svcCtx)
		err := l.ChangeAdminPassword(&req)
		httpy.ResultCtx(r, w, nil, err)
	}
}
