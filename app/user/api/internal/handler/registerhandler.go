package handler

import (
	"awesomeProject/common/xresp"
	"net/http"

	"awesomeProject/app/user/api/internal/logic"
	"awesomeProject/app/user/api/internal/svc"
	"awesomeProject/app/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		if err := xresp.Validate.StructCtx(r.Context(), req); err != nil {
			xresp.Response(r, w, nil, err)
			return
		}

		l := logic.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		xresp.Response(r, w, resp, err)
	}
}
