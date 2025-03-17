package logic_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"

	"pyth-go/app/pyth-common/types"
	"pyth-go/app/pyth-handler/rpc/internal/config"
	"pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/logic"
	"pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/service"
	dtypes "pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/types"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

func TestSlidingWindowLimitLogic(t *testing.T) {
	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	slidingWindowLogic := logic.NewSlidingWindowLimitLogic(svcCtx)

	taskInfo1 := &types.TaskInfo{
		MessageTemplateId: 1,
		Receiver:          []string{"receiver1-SW", "receiver2-SW"},
		ContentModel:      "test-contentModel-SW",
		SendChannel:       1,
		MsgType:           1,
	}
	cfg1 := dtypes.DeduplicationCfg{
		Num: 5,
	}
	contentDedupService := service.NewContentDeduplicationService(svcCtx, slidingWindowLogic)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < cfg1.Num; i++ {
		dupReceivers, err := slidingWindowLogic.Filter(ctx, contentDedupService, taskInfo1, cfg1)
		assert.Nil(t, err)
		t.Log(dupReceivers)
		assert.Equal(t, 0, len(dupReceivers))
	}
	dupReceivers, err := slidingWindowLogic.Filter(ctx, contentDedupService, taskInfo1, cfg1)

	for _ = range 2 {
		assert.Nil(t, err)
		t.Log(dupReceivers)
		assert.Equal(t, len(taskInfo1.Receiver), len(dupReceivers))
	}
}
