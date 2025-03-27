package task

import (
	"pyth-go/app/pyth-common/utils"
	"context"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"runtime"
)

type TaskPendingHolder struct {
	pool map[string]*ants.Pool
}

var defaultTaskPendingHolder = NewTaskPendingHolder()

func NewTaskPendingHolder() *TaskPendingHolder {
	// init all connection pool
	groupIds := utils.GetAllGroupIds()
	pool := make(map[string]*ants.Pool)
	size := runtime.NumCPU() * 2
	for _, value := range groupIds {
		var pushWorkerPool *ants.Pool
		if wp, err := ants.NewPool(size); err != nil {
			panic(fmt.Errorf("error occurred when creating push worker: %w", err))
		} else {
			pushWorkerPool = wp
		}
		pool[value] = pushWorkerPool
	}
	return &TaskPendingHolder{pool: pool}
}

// submit task to the corresponding pool
func (t TaskPendingHolder) Submit(ctx context.Context, groupId string, run TaskRun) error {
	return t.pool[groupId].Submit(func() {
		run.Run(ctx)
	})
}

func Submit(ctx context.Context, groupId string, run TaskRun) error {
	return defaultTaskPendingHolder.Submit(ctx, groupId, run)
}
