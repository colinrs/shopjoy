package warehouses

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/warehouses"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func SetDefaultWarehouseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetDefaultWarehouseReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := warehouses.NewSetDefaultWarehouseLogic(r.Context(), svcCtx)
		resp, err := l.SetDefaultWarehouse(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
