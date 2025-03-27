package handlers

import (
	"context"

	"pyth-go/app/pyth-common/types"
)

type ILimit interface {
	Limit(ctx context.Context, taskInfo types.TaskInfo) bool
}

type IHandler interface {
	ILimit
	Do(ctx context.Context, taskInfo types.TaskInfo) error
	Recall(ctx context.Context, taskInfo types.TaskInfo)
}


