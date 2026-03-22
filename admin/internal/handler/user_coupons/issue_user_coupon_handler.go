package user_coupons

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/user_coupons"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func IssueUserCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IssueUserCouponReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := user_coupons.NewIssueUserCouponLogic(r.Context(), svcCtx)
		resp, err := l.IssueUserCoupon(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
