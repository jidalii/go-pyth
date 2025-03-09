package app

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SendAccountModel = (*customSendAccountModel)(nil)

type (
	// SendAccountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSendAccountModel.
	SendAccountModel interface {
		sendAccountModel
		withSession(session sqlx.Session) SendAccountModel
	}

	customSendAccountModel struct {
		*defaultSendAccountModel
	}
)

// NewSendAccountModel returns a model for the database table.
func NewSendAccountModel(conn sqlx.SqlConn) SendAccountModel {
	return &customSendAccountModel{
		defaultSendAccountModel: newSendAccountModel(conn),
	}
}

func (m *customSendAccountModel) withSession(session sqlx.Session) SendAccountModel {
	return NewSendAccountModel(sqlx.NewSqlConnFromSession(session))
}
