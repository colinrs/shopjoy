package products

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/products"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func UpdateStockHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateStockReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := products.NewUpdateStockLogic(r.Context(), svcCtx)
		resp, err := l.UpdateStock(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
