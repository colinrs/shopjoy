package dashboard

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/dashboard"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func GetDashboardOverviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DashboardOverviewRequest
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := dashboard.NewGetDashboardOverviewLogic(r.Context(), svcCtx)
		resp, err := l.GetDashboardOverview(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
