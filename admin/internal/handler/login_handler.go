package handler

import (
	"net/http"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/httpy"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		tenantID, err := getTenantID(r)
		if err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		userResp, err := svcCtx.UserService.GetByEmail(r.Context(), tenantID, req.Email)
		if err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		accessToken, refreshToken, err := svcCtx.JWTManager.GenerateTokenPair(
			userResp.ID,
			shared.TenantID(userResp.TenantID),
			userResp.Email,
		)
		if err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}

		httpy.ResultCtx(r, w, types.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    3600,
			User: types.UserInfo{
				ID:     userResp.ID,
				Email:  userResp.Email,
				Name:   userResp.Name,
				Avatar: userResp.Avatar,
			},
		}, nil)
	}
}
