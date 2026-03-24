package webhooks

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/webhooks"
	"github.com/colinrs/shopjoy/admin/internal/svc"
)

func StripeWebhookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := webhooks.NewStripeWebhookLogic(r.Context(), svcCtx)
		err := l.StripeWebhook(r)
		httpy.ResultCtx(r, w, nil, err)
	}
}
