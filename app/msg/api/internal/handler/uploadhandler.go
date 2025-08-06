package handler

import (
	"awesomeProject/common/xresp"
	"net/http"

	"awesomeProject/app/msg/api/internal/logic"
	"awesomeProject/app/msg/api/internal/svc"
	"awesomeProject/app/msg/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func uploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		if err := xresp.Validate.StructCtx(r.Context(), req); err != nil {
			xresp.Response(r, w, nil, err)
			return
		}

		l := logic.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload(&req)
		xresp.Response(r, w, resp, err)
	}
}
