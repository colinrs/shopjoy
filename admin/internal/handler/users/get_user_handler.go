package users

import (
	"net/http"

	"github.com/colinrs/shopjoy/pkg/httpy"

	"github.com/colinrs/shopjoy/admin/internal/logic/users"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func GetUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserRequest
		if err := httpy.Parse(r, &req); err != nil {
			httpy.ResultCtx(r, w, nil, err)
			return
		}
		l := users.NewGetUserLogic(r.Context(), svcCtx)
		resp, err := l.GetUser(&req)
		httpy.ResultCtx(r, w, resp, err)
	}
}
