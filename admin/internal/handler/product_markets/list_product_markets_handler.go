// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product_markets

import (
	"net/http"

	"github.com/colinrs/shopjoy/admin/internal/logic/product_markets"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取商品市场配置列表
func ListProductMarketsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product_markets.NewListProductMarketsLogic(r.Context(), svcCtx)
		resp, err := l.ListProductMarkets()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
