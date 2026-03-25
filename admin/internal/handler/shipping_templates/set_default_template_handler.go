package shipping_templates

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/shipping_templates"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func SetDefaultTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetDefaultTemplateReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := shipping_templates.NewSetDefaultTemplateLogic(r.Context(), svcCtx)
		resp, err := l.SetDefaultTemplate(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
