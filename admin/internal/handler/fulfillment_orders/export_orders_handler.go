package fulfillment_orders

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/fulfillment_orders"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func ExportOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExportOrdersReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := fulfillment_orders.NewExportOrdersLogic(r.Context(), svcCtx, w, r)
		err := l.ExportOrders(&req)
		if err != nil {
			httpy.ResultCtx(r, w, nil, err)
		}
		// If successful, the CSV data has already been written to w
	}
}