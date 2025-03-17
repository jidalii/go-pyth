package services_test

import (
	"context"
	"flag"
	"strconv"
	"testing"
    "math/rand"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"

	"pyth-go/app/pyth-common/types"
	"pyth-go/app/pyth-handler/rpc/internal/config"
	"pyth-go/app/pyth-handler/rpc/internal/handler/services"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
	"pyth-go/common/zutils/arrayUtils"
)

var configFile = flag.String("f", "../../../etc/pythhandler_test.yaml", "the test config file")

const discardMessageKey = "discard_msg"

func addDiscardMessage(ctx context.Context, svcCtx *svc.ServiceContext, keys []string) error {
	_, err := svcCtx.RedisClient.SaddCtx(ctx, discardMessageKey, keys)
	return err
}

func delDiscardMessage(ctx context.Context, svcCtx *svc.ServiceContext, keys []string) error {
	_, err := svcCtx.RedisClient.SremCtx(ctx, discardMessageKey, keys)
	return err
}

func TestDiscardMessageService(t *testing.T) {
	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)

	ctx := context.Background()

	discardMessageService := services.NewDiscardMessageService(svcCtx)

	discardKeys := []string{"1", "2", "10", "19"}
	err := delDiscardMessage(ctx, svcCtx, discardKeys)
	assert.Nil(t, err)
	err = addDiscardMessage(ctx, svcCtx, discardKeys)
	assert.Nil(t, err)

	for i := 0; i < 20; i++ {
		taskInfo := &types.TaskInfo{
			MessageTemplateId: rand.Int63n(20),
		}

		res := discardMessageService.IsDiscard(ctx, taskInfo)
		if arrayUtils.ArrayIn(discardKeys, strconv.FormatInt(taskInfo.MessageTemplateId, 10)) {
            assert.True(t, res)
            t.Logf("discarded msg{id %d}", taskInfo.MessageTemplateId)
		} else {
			assert.False(t, res)
            t.Logf("not discarded msg{id %d}", taskInfo.MessageTemplateId)
		}
	}

    err = delDiscardMessage(ctx, svcCtx, discardKeys)
	assert.Nil(t, err)
}
