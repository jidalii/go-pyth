package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"pyth-go/app/pyth-handler/model"
	"pyth-go/app/pyth-handler/rpc/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis

	MessageTemplateModel model.MessageTemplateModel
	SendAccountModel     model.SendAccountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
		MessageTemplateModel: model.NewMessageTemplateModel(sqlConn, c.Cache),
		SendAccountModel: model.NewSendAccountModel(sqlConn, c.Cache),
	}
}
