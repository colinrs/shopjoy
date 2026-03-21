package inventory

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/inventory"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func BatchUpdateSafetyStockHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchUpdateSafetyStockReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := inventory.NewBatchUpdateSafetyStockLogic(r.Context(), svcCtx)
		resp, err := l.BatchUpdateSafetyStock(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
