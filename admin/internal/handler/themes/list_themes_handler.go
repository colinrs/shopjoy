package themes

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/themes"
	"github.com/colinrs/shopjoy/admin/internal/svc"
)

func ListThemesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := themes.NewListThemesLogic(r.Context(), svcCtx)
		resp, err := l.ListThemes()
		httpy.ResultCtx(r, w, resp, err)
	}
}
