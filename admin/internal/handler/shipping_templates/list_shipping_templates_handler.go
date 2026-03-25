package shipping_templates

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/shipping_templates"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func ListShippingTemplatesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListShippingTemplatesReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := shipping_templates.NewListShippingTemplatesLogic(r.Context(), svcCtx)
		resp, err := l.ListShippingTemplates(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
