package shop

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/shop"
	"github.com/colinrs/shopjoy/admin/internal/svc"
)

func GetBusinessHoursHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewGetBusinessHoursLogic(r.Context(), svcCtx)
		resp, err := l.GetBusinessHours()
		httpy.ResultCtx(r, w, resp, err)
	}
}
