// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product_markets

import (
	"net/http"

	"github.com/colinrs/shopjoy/admin/internal/logic/product_markets"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新商品市场配置
func UpdateProductMarketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateProductMarketReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := product_markets.NewUpdateProductMarketLogic(r.Context(), svcCtx)
		resp, err := l.UpdateProductMarket(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
