// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.6

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	messageTemplateFieldNames          = builder.RawFieldNames(&MessageTemplate{})
	messageTemplateRows                = strings.Join(messageTemplateFieldNames, ",")
	messageTemplateRowsExpectAutoSet   = strings.Join(stringx.Remove(messageTemplateFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	messageTemplateRowsWithPlaceHolder = strings.Join(stringx.Remove(messageTemplateFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheMessageTemplateIdPrefix = "cache:messageTemplate:id:"
)

type (
	messageTemplateModel interface {
		Insert(ctx context.Context, data *MessageTemplate) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*MessageTemplate, error)
		Update(ctx context.Context, data *MessageTemplate) error
		Delete(ctx context.Context, id int64) error
	}

	defaultMessageTemplateModel struct {
		sqlc.CachedConn
		table string
	}

	MessageTemplate struct {
		Id                  int64          `db:"id"`
		Name                string         `db:"name"`                 // 标题
		AuditStatus         int64          `db:"audit_status"`         // 当前消息审核状态： 10.待审核 20.审核成功 30.被拒绝
		FlowId              sql.NullString `db:"flow_id"`              // 工单ID
		MsgStatus           int64          `db:"msg_status"`           // 当前消息状态：10.新建 20.停用 30.启用 40.等待发送 50.发送中 60.发送成功 70.发送失败
		CronTaskId          sql.NullInt64  `db:"cron_task_id"`         // 定时任务Id (xxl-job-admin返回)
		CronCrowdPath       sql.NullString `db:"cron_crowd_path"`      // 定时发送人群的文件路径
		ExpectPushTime      sql.NullString `db:"expect_push_time"`     // 期望发送时间：0:立即发送 定时任务以及周期任务:cron表达式
		IdType              int64          `db:"id_type"`              // 消息的发送ID类型：10. userId 20.did 30.手机号 40.openId 50.email 60.企业微信userId
		SendChannel         int64          `db:"send_channel"`         // 消息发送渠道：10.IM 20.Push 30.短信 40.Email 50.公众号 60.小程序 70.企业微信
		TemplateType        int64          `db:"template_type"`        // 10.运营类 20.技术类接口调用
		MsgType             int64          `db:"msg_type"`             // 10.通知类消息 20.营销类消息 30.验证码类消息
		ShieldType          int64          `db:"shield_type"`          // 10.夜间不屏蔽 20.夜间屏蔽 30.夜间屏蔽(次日早上9点发送)
		MsgContent          string         `db:"msg_content"`          // 消息内容 占位符用{$var}表示
		SendAccount         int64          `db:"send_account"`         // 发送账号 一个渠道下可存在多个账号
		Creator             string         `db:"creator"`              // 创建者
		Updator             string         `db:"updator"`              // 更新者
		Auditor             string         `db:"auditor"`              // 审核人
		Team                string         `db:"team"`                 // 业务方团队
		Proposer            string         `db:"proposer"`             // 业务方
		IsDeleted           int64          `db:"is_deleted"`           // 是否删除：0.不删除 1.删除
		Created             int64          `db:"created"`              // 创建时间
		Updated             int64          `db:"updated"`              // 更新时间
		DeduplicationConfig string         `db:"deduplication_config"` // 限流配置
		TemplateSn          string         `db:"template_sn"`          // 发送消息的模版ID
	}
)

func newMessageTemplateModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultMessageTemplateModel {
	return &defaultMessageTemplateModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`message_template`",
	}
}

func (m *defaultMessageTemplateModel) Delete(ctx context.Context, id int64) error {
	messageTemplateIdKey := fmt.Sprintf("%s%v", cacheMessageTemplateIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, messageTemplateIdKey)
	return err
}

func (m *defaultMessageTemplateModel) FindOne(ctx context.Context, id int64) (*MessageTemplate, error) {
	messageTemplateIdKey := fmt.Sprintf("%s%v", cacheMessageTemplateIdPrefix, id)
	var resp MessageTemplate
	err := m.QueryRowCtx(ctx, &resp, messageTemplateIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", messageTemplateRows, m.table)
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

func (m *defaultMessageTemplateModel) Insert(ctx context.Context, data *MessageTemplate) (sql.Result, error) {
	messageTemplateIdKey := fmt.Sprintf("%s%v", cacheMessageTemplateIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, messageTemplateRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Name, data.AuditStatus, data.FlowId, data.MsgStatus, data.CronTaskId, data.CronCrowdPath, data.ExpectPushTime, data.IdType, data.SendChannel, data.TemplateType, data.MsgType, data.ShieldType, data.MsgContent, data.SendAccount, data.Creator, data.Updator, data.Auditor, data.Team, data.Proposer, data.IsDeleted, data.Created, data.Updated, data.DeduplicationConfig, data.TemplateSn)
	}, messageTemplateIdKey)
	return ret, err
}

func (m *defaultMessageTemplateModel) Update(ctx context.Context, data *MessageTemplate) error {
	messageTemplateIdKey := fmt.Sprintf("%s%v", cacheMessageTemplateIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, messageTemplateRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Name, data.AuditStatus, data.FlowId, data.MsgStatus, data.CronTaskId, data.CronCrowdPath, data.ExpectPushTime, data.IdType, data.SendChannel, data.TemplateType, data.MsgType, data.ShieldType, data.MsgContent, data.SendAccount, data.Creator, data.Updator, data.Auditor, data.Team, data.Proposer, data.IsDeleted, data.Created, data.Updated, data.DeduplicationConfig, data.TemplateSn, data.Id)
	}, messageTemplateIdKey)
	return err
}

func (m *defaultMessageTemplateModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheMessageTemplateIdPrefix, primary)
}

func (m *defaultMessageTemplateModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", messageTemplateRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultMessageTemplateModel) tableName() string {
	return m.table
}
