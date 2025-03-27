package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"

	"pyth-go/app/pyth-support/mq/consumer"
)

type Config struct {
	zrpc.RpcServerConf
	Cache cache.CacheConf
	DB    struct {
		DataSource string
	}
	Consumer consumer.ConsumerGroupConfig
}
