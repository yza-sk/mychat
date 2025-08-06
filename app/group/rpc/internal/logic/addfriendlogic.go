package logic

import (
	"awesomeProject/app/group/model"
	"awesomeProject/common/biz"
	"awesomeProject/common/xerr"
	"awesomeProject/common/xmq"
	"awesomeProject/proto/group"
	"awesomeProject/proto/msg"
	"context"
	"github.com/pkg/errors"
	"strings"
	"time"

	"awesomeProject/app/group/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendLogic {
	return &AddFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送添加好友请求
func (l *AddFriendLogic) AddFriend(in *group.AddFriendRequest) (*group.AddFriendResponse, error) {
	fromUid := in.FromUid
	toUid := in.ToUid
	// 生成groupId
	groupId := biz.GetGroupId(fromUid, toUid)
	// 查询这两个用户的nickName
	u1, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(fromUid))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "AddFriend query user failed, fromUid: %v", fromUid)
	}
	u2, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(toUid))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "AddFriend query user failed, toUid: %v", toUid)
	}
	groupName := strings.Join([]string{u1.NickName, u2.NickName}, ", ")
	// 创建一个group
	_, err = l.svcCtx.GroupModel.Insert(l.ctx, &model.Group{
		Id:     groupId,
		Name:   groupName,
		Type:   1,
		Status: model.GroupStatusNo,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "AddFriend insert group failed: %v", err)
	}
	// 加好友请求消息 放入消息队列
	chatMsg := &msg.ChatMsg{
		GroupId:    groupId,
		SenderId:   fromUid,
		Type:       0,
		Content:    "请求加你为好友",
		Uuid:       biz.GetUuid(),
		CreateTime: time.Now().UnixMilli(),
	}
	err = xmq.PushToMq(l.ctx, l.svcCtx.MqWriter, chatMsg)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "AddFriend push to mq failed: %v", err)
	}
	return &group.AddFriendResponse{
		GroupId: groupId,
	}, nil
}
