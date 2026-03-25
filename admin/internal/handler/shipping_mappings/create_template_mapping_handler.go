package shipping_mappings

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/shipping_mappings"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func CreateTemplateMappingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateTemplateMappingReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := shipping_mappings.NewCreateTemplateMappingLogic(r.Context(), svcCtx)
		resp, err := l.CreateTemplateMapping(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
