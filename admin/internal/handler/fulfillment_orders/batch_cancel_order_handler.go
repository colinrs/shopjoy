package fulfillment_orders

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/fulfillment_orders"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func BatchCancelOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchCancelOrderReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := fulfillment_orders.NewBatchCancelOrderLogic(r.Context(), svcCtx)
		resp, err := l.BatchCancelOrder(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
