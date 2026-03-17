// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/colinrs/shopjoy/shop/internal/logic"
	"github.com/colinrs/shopjoy/shop/internal/svc"
	"github.com/colinrs/shopjoy/shop/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新库存
func UpdateStockHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateStockReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUpdateStockLogic(r.Context(), svcCtx)
		resp, err := l.UpdateStock(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
