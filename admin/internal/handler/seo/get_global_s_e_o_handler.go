package seo

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/seo"
	"github.com/colinrs/shopjoy/admin/internal/svc"
)

func GetGlobalSEOHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seo.NewGetGlobalSEOLogic(r.Context(), svcCtx)
		resp, err := l.GetGlobalSEO()
		httpy.ResultCtx(r, w, resp, err)
	}
}
