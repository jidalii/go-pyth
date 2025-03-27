package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SendAccountModel = (*customSendAccountModel)(nil)

type (
	// SendAccountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSendAccountModel.
	SendAccountModel interface {
		sendAccountModel
	}

	customSendAccountModel struct {
		*defaultSendAccountModel
	}
)

// NewSendAccountModel returns a model for the database table.
func NewSendAccountModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SendAccountModel {
	return &customSendAccountModel{
		defaultSendAccountModel: newSendAccountModel(conn, c, opts...),
	}
}
