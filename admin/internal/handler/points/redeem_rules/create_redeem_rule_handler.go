package redeem_rules

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/points/redeem_rules"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func CreateRedeemRuleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateRedeemRuleReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := redeem_rules.NewCreateRedeemRuleLogic(r.Context(), svcCtx)
		resp, err := l.CreateRedeemRule(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
