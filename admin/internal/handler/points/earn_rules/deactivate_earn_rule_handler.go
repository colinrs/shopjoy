package earn_rules

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/points/earn_rules"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func DeactivateEarnRuleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeactivateEarnRuleReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := earn_rules.NewDeactivateEarnRuleLogic(r.Context(), svcCtx)
		resp, err := l.DeactivateEarnRule(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
