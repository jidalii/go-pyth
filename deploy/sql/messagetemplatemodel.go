package app

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ MessageTemplateModel = (*customMessageTemplateModel)(nil)

type (
	// MessageTemplateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageTemplateModel.
	MessageTemplateModel interface {
		messageTemplateModel
		withSession(session sqlx.Session) MessageTemplateModel
	}

	customMessageTemplateModel struct {
		*defaultMessageTemplateModel
	}
)

// NewMessageTemplateModel returns a model for the database table.
func NewMessageTemplateModel(conn sqlx.SqlConn) MessageTemplateModel {
	return &customMessageTemplateModel{
		defaultMessageTemplateModel: newMessageTemplateModel(conn),
	}
}

func (m *customMessageTemplateModel) withSession(session sqlx.Session) MessageTemplateModel {
	return NewMessageTemplateModel(sqlx.NewSqlConnFromSession(session))
}
