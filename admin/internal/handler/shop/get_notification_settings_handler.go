package shop

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/shop"
	"github.com/colinrs/shopjoy/admin/internal/svc"
)

func GetNotificationSettingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := shop.NewGetNotificationSettingsLogic(r.Context(), svcCtx)
		resp, err := l.GetNotificationSettings()
		httpy.ResultCtx(r, w, resp, err)
	}
}
