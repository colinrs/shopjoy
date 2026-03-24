package users

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/users"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func ExportUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExportUsersRequest
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := users.NewExportUsersLogic(r.Context(), svcCtx, w, r)
		err := l.ExportUsers(&req)
		if err != nil {
			httpy.ResultCtx(r, w, nil, err)
		}
		// If successful, the CSV data has already been written to w
	}
}
