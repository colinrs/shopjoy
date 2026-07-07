// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package products

import (
	"net/http"

	"github.com/colinrs/shopjoy/admin/internal/logic/products"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 搜索SKU（下拉选择用）
func SearchSKUsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchSKUsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := products.NewSearchSKUsLogic(r.Context(), svcCtx)
		resp, err := l.SearchSKUs(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
