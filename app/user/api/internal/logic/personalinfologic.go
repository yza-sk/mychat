package logic

import (
	"awesomeProject/app/user/rpc/userclient"
	"awesomeProject/common/ctxdata"
	"context"
	"github.com/jinzhu/copier"

	"awesomeProject/app/user/api/internal/svc"
	"awesomeProject/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PersonalInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPersonalInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PersonalInfoLogic {
	return &PersonalInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PersonalInfoLogic) PersonalInfo(req *types.PersonalInfoRequest) (*types.PersonalInfoResponse, error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	personInfo, err := l.svcCtx.UserRpc.PersonalInfo(l.ctx, &userclient.PersonalInfoRequest{
		Id: uid,
	})
	if err != nil {
		return nil, err
	}
	var resp types.PersonalInfoResponse
	copier.Copy(&resp, personInfo)
	return &resp, nil
}
