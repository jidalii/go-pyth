package task

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"pyth-go/app/pyth-common/types"
	"pyth-go/app/pyth-handler/rpc/internal/handler/handlers"
	"pyth-go/app/pyth-handler/rpc/internal/handler/services"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

type TaskRun interface {
	Run(ctx context.Context)
}

type Task struct {
	TaskInfo types.TaskInfo
	svcCtx   *svc.ServiceContext
}

func NewTask(taskInfo types.TaskInfo, svcCtx *svc.ServiceContext) *Task {
	return &Task{TaskInfo: taskInfo, svcCtx: svcCtx}
}

func (t Task) Run(ctx context.Context) {

	if services.NewDiscardMessageService(t.svcCtx).IsDiscard(ctx, &t.TaskInfo) {
		logx.WithContext(ctx).Infow("Taks is discarded", logx.Field("task_info", t.TaskInfo))
		return
	}
	// 1. Shield tasks
	services.NewShieldService(t.svcCtx).Shield(ctx, &t.TaskInfo, time.Now())
	// 2. Deduplication 1. same content in N minutes, 2. N times for the same channel in a day
	if len(t.TaskInfo.Receiver) > 0 {
		dedupSvc, err := services.NewDeduplicationService(t.svcCtx)
		if err != nil {
			logx.Errorw("NewDeduplicationService err", logx.Field("err", err))
			return
		}
		dedupSvc.Do(ctx, &t.TaskInfo)
	}
	// 3. Do the task
	if len(t.TaskInfo.Receiver) > 0 {
		h := handlers.GetHandler(t.TaskInfo.SendChannel)
		for {
			if h.Limit(ctx, t.TaskInfo) {
				err := h.Do(ctx, t.TaskInfo)
				if err != nil {
					logx.Errorw("DoHandler err", logx.Field("task_info", t.TaskInfo), logx.Field("err", err))
				}
				return
			}
		}
	}
}
