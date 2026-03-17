package handler

import (
	"net/http"

	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/httpy"
)

func ChangePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangePasswordRequest
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		tenantID, err := getTenantID(r)
		if err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		var userID int64 = 1

		appReq := appUser.ChangePasswordRequest{
			UserID:          userID,
			TenantID:        tenantID,
			OldPassword:     req.OldPassword,
			NewPassword:     req.NewPassword,
			ConfirmPassword: req.ConfirmPassword,
		}

		err = svcCtx.UserService.ChangePassword(r.Context(), appReq)
		httpy.ResultCtx(r, w, nil, err)
	}
}
