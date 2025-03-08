package types

import (
	"context"

	"pyth-go/app/pyth-common/types"
)

type IDeduplicationService interface {
	DoDeduplication(ctx context.Context, taskInfo *types.TaskInfo, cfg DeduplicationCfg) error
	GetDeduplicationSingleKey(taskInfo *types.TaskInfo, receiver string) string
}
