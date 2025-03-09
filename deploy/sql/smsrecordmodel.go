package app

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SmsRecordModel = (*customSmsRecordModel)(nil)

type (
	// SmsRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSmsRecordModel.
	SmsRecordModel interface {
		smsRecordModel
		withSession(session sqlx.Session) SmsRecordModel
	}

	customSmsRecordModel struct {
		*defaultSmsRecordModel
	}
)

// NewSmsRecordModel returns a model for the database table.
func NewSmsRecordModel(conn sqlx.SqlConn) SmsRecordModel {
	return &customSmsRecordModel{
		defaultSmsRecordModel: newSmsRecordModel(conn),
	}
}

func (m *customSmsRecordModel) withSession(session sqlx.Session) SmsRecordModel {
	return NewSmsRecordModel(sqlx.NewSqlConnFromSession(session))
}
