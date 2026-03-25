package dashboard

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/dashboard"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func GetRecentActivitiesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecentActivitiesRequest
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := dashboard.NewGetRecentActivitiesLogic(r.Context(), svcCtx)
		resp, err := l.GetRecentActivities(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
