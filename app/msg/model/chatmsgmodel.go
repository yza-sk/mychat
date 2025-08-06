package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatMsgModel = (*customChatMsgModel)(nil)

type (
	// ChatMsgModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatMsgModel.
	ChatMsgModel interface {
		chatMsgModel
		FindMsgListByLastMsgId(ctx context.Context, groupId string, minMsgId int64, maxMsgId int64) ([]*ChatMsg, error)
		FindLastMsgByGroupId(ctx context.Context, groupId string) (*ChatMsg, error)
		// 新增事务方法 [新增位置1]
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		TransInsert(ctx context.Context, session sqlx.Session, data *ChatMsg) (sql.Result, error)
		TransFindOne(ctx context.Context, session sqlx.Session, id int64) (*ChatMsg, error)
	}

	customChatMsgModel struct {
		*defaultChatMsgModel
	}
)

// NewChatMsgModel returns a model for the database table.
func NewChatMsgModel(conn sqlx.SqlConn, c cache.CacheConf) ChatMsgModel {
	return &customChatMsgModel{
		defaultChatMsgModel: newChatMsgModel(conn, c),
	}
}

// 获取指定群组的离线消息列表
func (m *customChatMsgModel) FindMsgListByLastMsgId(ctx context.Context, groupId string, minMsgId int64, maxMsgId int64) ([]*ChatMsg, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ? and `id` > ? and `id` < ? order by `id` desc limit %d", chatMsgRows, m.table, PerPullNum)
	var resp []*ChatMsg
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, groupId, minMsgId, maxMsgId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 获取指定群组的 最后一条消息
func (m *customChatMsgModel) FindLastMsgByGroupId(ctx context.Context, groupId string) (*ChatMsg, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ? order by `id` desc limit 1", chatMsgRows, m.table)
	var resp ChatMsg
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, groupId)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

// 新增事务执行方法 [新增位置2]
func (m *defaultChatMsgModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

// 新增事务内插入方法 [新增位置3]
func (m *defaultChatMsgModel) TransInsert(ctx context.Context, session sqlx.Session, data *ChatMsg) (sql.Result, error) {
	chatMsgIdKey := fmt.Sprintf("%s%v", cacheChatMsgIdPrefix, data.Id)
	chatMsgUuidKey := fmt.Sprintf("%s%v", cacheChatMsgUuidPrefix, data.Uuid)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, chatMsgRowsExpectAutoSet)
		if session != nil {
			return session.ExecCtx(ctx, query, data.GroupId, data.SenderId, data.Type, data.Content, data.Uuid)
		}
		return conn.ExecCtx(ctx, query, data.GroupId, data.SenderId, data.Type, data.Content, data.Uuid)
	}, chatMsgIdKey, chatMsgUuidKey)
}

func (m *defaultChatMsgModel) TransFindOne(ctx context.Context, session sqlx.Session, id int64) (*ChatMsg, error) {
	chatMsgIdKey := fmt.Sprintf("%s%v", cacheChatMsgIdPrefix, id)
	var resp ChatMsg
	err := m.QueryRowCtx(ctx, &resp, chatMsgIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", chatMsgRows, m.table)
		if session != nil {
			return session.QueryRowCtx(ctx, v, query, id)
		}
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
