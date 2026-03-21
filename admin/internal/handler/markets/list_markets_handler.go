package markets

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/markets"
	"github.com/colinrs/shopjoy/admin/internal/svc"
)

func ListMarketsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := markets.NewListMarketsLogic(r.Context(), svcCtx)
		resp, err := l.ListMarkets()
		httpy.ResultCtx(r, w, resp, err)
	}
}
