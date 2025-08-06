package logic

import (
	"awesomeProject/common/ctxdata"
	"awesomeProject/proto/msg"
	"context"
	"github.com/jinzhu/copier"

	"awesomeProject/app/msg/api/internal/svc"
	"awesomeProject/app/msg/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPullLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullLogic {
	return &PullLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PullLogic) Pull(req *types.PullRequest) (*types.PullResponse, error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	var pbPullRequest msg.PullRequest
	copier.Copy(&pbPullRequest, req)
	pbPullRequest.UserId = uid
	pbPullResponse, err := l.svcCtx.MsgRpc.Pull(l.ctx, &pbPullRequest)
	if err != nil {
		return nil, err
	}
	var resp types.PullResponse
	err = copier.Copy(&resp, pbPullResponse)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
