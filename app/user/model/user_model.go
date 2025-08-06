package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		// 添加事务方法
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		TransInsert(ctx context.Context, session sqlx.Session, data *User) (sql.Result, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

// 实现事务执行方法
func (m *customUserModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

// 实现事务内插入方法
func (m *customUserModel) TransInsert(ctx context.Context, session sqlx.Session, data *User) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, userRowsExpectAutoSet)
	return session.ExecCtx(ctx, query, data.Email, data.Password, data.NickName, data.Gender, data.AvatarUrl)
}
