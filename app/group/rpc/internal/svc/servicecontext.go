package svc

import (
	modelGroup "awesomeProject/app/group/model"
	"awesomeProject/app/group/rpc/internal/config"
	modelMsg "awesomeProject/app/msg/model"
	modelUser "awesomeProject/app/user/model"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

type ServiceContext struct {
	Config         config.Config
	GroupModel     modelGroup.GroupModel
	GroupUserModel modelGroup.GroupUserModel
	UserModel      modelUser.UserModel
	ChatMsgModel   modelMsg.ChatMsgModel
	MqWriter       *kafka.Writer
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Db.DataSource)
	mqWriter := &kafka.Writer{
		Addr:         kafka.TCP(c.MqConf.Brokers...),
		Topic:        c.MqConf.Topic,
		BatchTimeout: time.Millisecond * 20,
	}
	return &ServiceContext{
		Config:         c,
		GroupModel:     modelGroup.NewGroupModel(sqlConn, c.Cache),
		GroupUserModel: modelGroup.NewGroupUserModel(sqlConn, c.Cache),
		UserModel:      modelUser.NewUserModel(sqlConn, c.Cache),
		ChatMsgModel:   modelMsg.NewChatMsgModel(sqlConn, c.Cache),
		MqWriter:       mqWriter,
	}
}
