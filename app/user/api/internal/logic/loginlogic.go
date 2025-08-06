package logic

import (
	"awesomeProject/app/user/rpc/userclient"
	"context"
	"github.com/jinzhu/copier"
	"awesomeProject/app/user/api/internal/svc"
	"awesomeProject/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	loginResp, err := l.svcCtx.UserRpc.Login(l.ctx, &userclient.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	var resp types.LoginResponse
	copier.Copy(&resp, loginResp)
	return &resp, nil
}
