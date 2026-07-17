// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user_coupons

import (
	"net/http"

	"github.com/colinrs/shopjoy/admin/internal/logic/user_coupons"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 批量发放优惠券
func BatchIssueUserCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchIssueUserCouponReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user_coupons.NewBatchIssueUserCouponLogic(r.Context(), svcCtx)
		resp, err := l.BatchIssueUserCoupon(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
