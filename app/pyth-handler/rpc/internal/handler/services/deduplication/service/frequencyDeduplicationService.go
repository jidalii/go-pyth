package service

import (
	"context"
    "fmt"


	"pyth-go/app/pyth-common/types"
	"pyth-go/common/zutils/arrayUtils"
	dtypes "pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/types"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

const frequencyDedupPrefix = "FRE"

// A frequency deduplication service
//
// Each day, a user can only receive messages of a channel for at most N times
//
// Use slidingWindowLimitLogic to filter the receivers
type frequencyDeduplicationService struct {
    svcCtx *svc.ServiceContext
	limit  dtypes.LimitLogic
}

func NewFrequencyDeduplicationService(svcCtx *svc.ServiceContext, limit dtypes.LimitLogic) dtypes.IDeduplicationService {
	return &frequencyDeduplicationService{
		svcCtx: svcCtx,
		limit:  limit,
	}
}

func (s frequencyDeduplicationService) DoDeduplication(ctx context.Context, taskInfo *types.TaskInfo, cfg dtypes.DeduplicationCfg) error {
	var newRows []string
	filter, err := s.limit.Filter(ctx, s, taskInfo, cfg)
	if err != nil {
		return err
	}
	for _, s := range taskInfo.Receiver {
		if !arrayUtils.ArrayIn(filter, s) {
			newRows = append(newRows, s)
		}
	}
	taskInfo.Receiver = newRows
	return nil
}

func (s frequencyDeduplicationService) GetDeduplicationSingleKey(taskInfo *types.TaskInfo, receiver string) string {
	return fmt.Sprintf("%s_%s_%d_%d", frequencyDedupPrefix, receiver, taskInfo.MessageTemplateId, taskInfo.SendChannel)
}
