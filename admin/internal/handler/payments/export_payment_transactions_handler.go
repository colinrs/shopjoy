package payments

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/payments"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func ExportPaymentTransactionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExportPaymentTransactionsReq
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := payments.NewExportPaymentTransactionsLogic(r.Context(), svcCtx, w, r)
		err := l.ExportPaymentTransactions(&req)
		if err != nil {
			httpy.ResultCtx(r, w, nil, err)
		}
		// If successful, the CSV data has already been written to w
	}
}
