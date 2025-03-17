package services

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/jsonx"
	// "github.com/spf13/cast"
	"github.com/panjf2000/ants/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"pyth-go/app/pyth-common/types"
	// "pyth-go/common/encrypt"
	// "pyth-go/common/zutils/arrayUtils"
	"pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/logic"
	"pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/service"
	dtypes "pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/types"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

type DeduplicationService struct {
	svcCtx     *svc.ServiceContext
	workerPool *ants.Pool
	services   map[string]dtypes.IDeduplicationService
	done       chan struct{}
}

const (
	Content        = "10"
	Frequency      = "20"
	dedupKeyPrefix = "deduplication"
	defaultTimeout = 30 * time.Second
)

func NewDeduplicationService(svcCtx *svc.ServiceContext) (*DeduplicationService, error) {
	pool, err := ants.NewPool(
		10,
		ants.WithNonblocking(true),
		ants.WithExpiryDuration(time.Minute),
	)
	if err != nil {
		logx.Errorw("failed to create ants pool", logx.Field("err", err))
	}

	services := make(map[string]dtypes.IDeduplicationService, 2)
	services[fmt.Sprintf("%s_%s", dedupKeyPrefix, Content)] =
		service.NewContentDeduplicationService(svcCtx, logic.NewSlidingWindowLimitLogic(svcCtx))
	services[fmt.Sprintf("%s_%s", dedupKeyPrefix, Frequency)] =
		service.NewFrequencyDeduplicationService(svcCtx, logic.NewTokenBucketLimitLogic(svcCtx))

	svc := &DeduplicationService{
		svcCtx:     svcCtx,
		workerPool: pool,
		services:   services,
		done:       make(chan struct{}),
	}

	return svc, nil
}

func (h DeduplicationService) Do(ctx context.Context, taskInfo *types.TaskInfo) {
	temp, err := h.svcCtx.MessageTemplateModel.FindOne(ctx, taskInfo.MessageTemplateId)
	if err != nil {
		logx.Errorw("deduplicationRuleService: failed to find template", logx.Field("err", err))
		return
	}
	if temp.DeduplicationConfig == "" {
		// empty deduplication config, return directly
		return
	}

	dedupConfig := make(map[string]dtypes.DeduplicationCfg)
	err = jsonx.Unmarshal([]byte(temp.DeduplicationConfig), &dedupConfig)
	if err != nil {
		logx.Errorw("deduplicationRuleService jsonx.Unmarshal err", logx.Field("err", err))
		return
	}

	if len(dedupConfig) <= 0 {
		return
	}


	for key, cfg := range dedupConfig {
		if err := h.execDedup(ctx, key, taskInfo, cfg); err != nil {
			logx.Errorw("deduplicationRuleService: failed to exec deduplication", logx.Field("err", err))
		}
	}
}

func (h DeduplicationService) Close() {
	close(h.done)
	h.workerPool.Release()
}

func (h DeduplicationService) execDedup(ctx context.Context, key string, taskInfo *types.TaskInfo, cfg dtypes.DeduplicationCfg) error {
	service, ok := h.services[key]
	if !ok {
		logx.Errorw("deduplicationRuleService: unknown deduplication key", logx.Field("key", key))
		return fmt.Errorf("unknown deduplication key: %s", key)
	}
	if err := service.DoDeduplication(ctx, taskInfo, cfg); err != nil {
		logx.Errorw("deduplication failed", logx.Field("err", err))
	}

    return nil
}
