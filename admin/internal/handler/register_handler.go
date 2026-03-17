package handler

import (
	"net/http"

	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/httpy"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterRequest
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		tenantID, err := getTenantID(r)
		if err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		appReq := appUser.CreateUserRequest{
			TenantID: tenantID,
			Email:    req.Email,
			Phone:    req.Phone,
			Password: req.Password,
			Name:     req.Name,
		}

		resp, err := svcCtx.UserService.Register(r.Context(), appReq)
		if err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		httpy.ResultCtx(r, w, types.RegisterResponse{
			ID:        resp.ID,
			Email:     resp.Email,
			Name:      resp.Name,
			CreatedAt: resp.CreatedAt,
		}, nil)
	}
}
