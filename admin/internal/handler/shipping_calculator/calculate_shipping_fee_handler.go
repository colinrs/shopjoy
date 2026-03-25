package shipping_calculator

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/shipping_calculator"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func CalculateShippingFeeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CalculateShippingFeeReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := shipping_calculator.NewCalculateShippingFeeLogic(r.Context(), svcCtx)
		resp, err := l.CalculateShippingFee(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
