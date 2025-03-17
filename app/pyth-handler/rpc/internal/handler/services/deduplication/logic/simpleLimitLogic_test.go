package logic_test

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/stretchr/testify/assert"

	"pyth-go/app/pyth-common/types"
	"pyth-go/app/pyth-handler/rpc/internal/config"
	"pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/logic"
	"pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/service"
	dtypes "pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/types"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

var configFile = flag.String("f", "../../../../../etc/pythhandler_test.yaml", "the test config file")

func TestSimpleLimitLogic(t *testing.T) {
	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	simpleLimitLogic := logic.NewSimpleLimitLogic(svcCtx)

    taskInfo1 := &types.TaskInfo{
		MessageTemplateId: 1,
        Receiver: []string{"receiver1-SP", "receiver2-SP"},
		ContentModel: "test-contentModel-SP",
		SendChannel: 1,
		MsgType: 1,
    }
    cfg1 := dtypes.DeduplicationCfg{
        Num: 3,
        Time: 20,
    }
    contentDedupService := service.NewContentDeduplicationService(svcCtx, simpleLimitLogic)

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	for _ = range cfg1.Num {
		dupReceivers, err := simpleLimitLogic.Filter(ctx, contentDedupService, taskInfo1, cfg1)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(dupReceivers))
	}
	dupReceivers, err := simpleLimitLogic.Filter(ctx, contentDedupService, taskInfo1, cfg1)
	assert.Nil(t, err)
	assert.Equal(t, len(taskInfo1.Receiver), len(dupReceivers))
}