package logic

import (
	"awesomeProject/app/msg/model"
	"awesomeProject/common/xerr"
	"awesomeProject/common/xmq"
	"awesomeProject/proto/msg"
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"awesomeProject/app/msg/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadLogic) Upload(in *msg.UploadRequest) (*msg.UploadResponse, error) {
	// 给消息结构体标准化
	chatMsg := &model.ChatMsg{
		GroupId:  in.GroupId,
		SenderId: in.SenderId,
		Type:     in.Type,
		Content:  in.Content,
		Uuid:     in.Uuid,
	}

	err := l.svcCtx.ChatMsgModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 将消息体存到数据库
		ret, err := l.svcCtx.ChatMsgModel.TransInsert(l.ctx, session, chatMsg)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCodeMsg(xerr.DB_ERROR, "消息uuid已存在"), "insert message failed, msg: %+v", chatMsg)
		}
		// 获取消息体存到数据库中的主键id
		val, err := ret.LastInsertId()
		chatMsg.Id = uint64(val)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCodeMsg(xerr.DB_ERROR, "获取消息id失败"), "get message id failed, msg: %+v", chatMsg)
		}
		// 验证消息插入数据库成功并获取完整的消息对象
		dbMsg, err := l.svcCtx.ChatMsgModel.TransFindOne(l.ctx, session, int64(chatMsg.Id))
		if err != nil {
			return errors.Wrapf(xerr.NewErrCodeMsg(xerr.DB_ERROR, "获取消息失败"), "get message failed, msg: %+v", chatMsg)
		}
		// 放入消息队列
		var mqMsg msg.ChatMsg
		// 将chatmsg类型转换mqmsg类型
		copier.Copy(&mqMsg, dbMsg)
		mqMsg.CreateTime = dbMsg.CreateTime.UnixMilli()
		chatMsg.CreateTime = dbMsg.CreateTime
		// 推送进入消息队列
		err = xmq.PushToMq(l.ctx, l.svcCtx.MqWriter, &mqMsg)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCodeMsg(xerr.MQ_ERROR, "消息推送失败"), "push message to mq failed, msg: %+v, err: %v", chatMsg, err)
		}
		logx.Infof("push to mq msg: %+v", chatMsg)
		// commit
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &msg.UploadResponse{
		Id:         int64(chatMsg.Id),
		CreateTime: chatMsg.CreateTime.UnixMilli(),
	}, nil
}
