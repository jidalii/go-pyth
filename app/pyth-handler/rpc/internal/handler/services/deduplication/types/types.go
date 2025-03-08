package types

import (
	"context"

	"pyth-go/app/pyth-common/types"
)

type DeduplicationCfg struct {
	Num  int   `json:"num"`  // The number of messages allowed to be sent in the time window
	Time int64 `json:"time"` // The time window for deduplication (in seconds)
}

type DeduplicationService interface {
	Deduplication(ctx context.Context, taskInfo *types.TaskInfo, cfg DeduplicationCfg) error
}

type LimitLogic interface {
	Filter(ctx context.Context, service IDeduplicationService, taskInfo *types.TaskInfo, cfg DeduplicationCfg) (dupReceivers []string, err error)
}
