package fulfillment_statistics

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/fulfillment_statistics"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func GetFulfillmentStatisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRefundStatisticsReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := fulfillment_statistics.NewGetFulfillmentStatisticsLogic(r.Context(), svcCtx)
		resp, err := l.GetFulfillmentStatistics(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
