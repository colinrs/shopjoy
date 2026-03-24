package themes

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/themes"
	"github.com/colinrs/shopjoy/admin/internal/svc"
)

func GetCurrentThemeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := themes.NewGetCurrentThemeLogic(r.Context(), svcCtx)
		resp, err := l.GetCurrentTheme()
		httpy.ResultCtx(r, w, resp, err)
	}
}
