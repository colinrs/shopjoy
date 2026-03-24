package pages

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/pages"
	"github.com/colinrs/shopjoy/admin/internal/svc"
)

func ListPagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := pages.NewListPagesLogic(r.Context(), svcCtx)
		resp, err := l.ListPages()
		httpy.ResultCtx(r, w, resp, err)
	}
}
