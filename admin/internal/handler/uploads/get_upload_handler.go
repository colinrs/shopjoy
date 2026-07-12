// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package uploads

import (
	"net/http"

	"github.com/colinrs/shopjoy/admin/internal/logic/uploads"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取图片元数据
func GetUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := uploads.NewGetUploadLogic(r.Context(), svcCtx)
		resp, err := l.GetUpload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
