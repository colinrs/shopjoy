package seo

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/seo"
	"github.com/colinrs/shopjoy/admin/internal/svc"
)

func ListPageSEOHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seo.NewListPageSEOLogic(r.Context(), svcCtx)
		resp, err := l.ListPageSEO()
		httpy.ResultCtx(r, w, resp, err)
	}
}
