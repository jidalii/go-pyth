package service

import (
    "context"

    "pyth-go/app/pyth-common/types"
    "pyth-go/common/zutils/arrayUtils"

    dtypes "pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/types"
)

type deduplicationService struct {
}

func (c deduplicationService) Deduplication(ctx context.Context,
	limit dtypes.LimitLogic,
	service dtypes.IDeduplicationService,
	taskInfo *types.TaskInfo,
	cfg dtypes.DeduplicationCfg) error {

	var newRows []string
	filter, err := limit.Filter(ctx, service, taskInfo, cfg)
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