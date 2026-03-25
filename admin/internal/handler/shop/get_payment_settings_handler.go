package shop

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/shop"
	"github.com/colinrs/shopjoy/admin/internal/svc"
)

func GetPaymentSettingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewGetPaymentSettingsLogic(r.Context(), svcCtx)
		resp, err := l.GetPaymentSettings()
		httpy.ResultCtx(r, w, resp, err)
	}
}
