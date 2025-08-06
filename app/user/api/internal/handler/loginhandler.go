package handler

import (
	"awesomeProject/common/xresp"
	"net/http"

	"awesomeProject/app/user/api/internal/logic"
	"awesomeProject/app/user/api/internal/svc"
	"awesomeProject/app/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 对请求体结构参数进行规则性校验，若必要字段为空则返回错误
		if err := xresp.Validate.StructCtx(r.Context(), req); err != nil {
			xresp.Response(r, w, nil, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		xresp.Response(r, w, resp, err)
	}
}
