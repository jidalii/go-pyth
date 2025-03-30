package handler

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"

	"pyth-go/app/pyth-common/enums/chanType"
	"pyth-go/app/pyth-common/enums/msgType"
	"pyth-go/app/pyth-common/types"
	"pyth-go/app/pyth-handler/rpc/internal/handler/task"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
	"pyth-go/app/pyth-support/mq/consumer"
)

type TaskExecutor struct {
	svcCtx        *svc.ServiceContext
	ConsumerGroup *consumer.MQConsumerGroup
	Consumer      *consumer.MQConsumer
	isGroup       bool
}

func NewTaskExecutor(svcCtx *svc.ServiceContext, isGroup bool) *TaskExecutor {
	executor := &TaskExecutor{
		svcCtx:  svcCtx,
		isGroup: isGroup,
	}
	if isGroup {
		consumerGroup := consumer.NewConsumerGroup(svcCtx.Config.Consumer)
		executor.ConsumerGroup = &consumerGroup
	} else {
		consumer := consumer.NewConsumer(svcCtx.Config.Consumer)
		executor.Consumer = &consumer
	}
	return executor
}

func (t *TaskExecutor) Start() {
    ctx := context.Background()
	if t.isGroup {
		t.ConsumerGroup.Start(ctx, t.OnMessage)
	} else {
		go t.Consumer.Start(ctx, t.OnMessage)
	}
}

func (t *TaskExecutor) OnMessage(m kafka.Message) error {
	ctx := context.Background()
	var taskList []types.TaskInfo
	_ = jsonx.Unmarshal(m.Value, &taskList)
	for _, taskInfo := range taskList {
		logx.WithContext(ctx).Infow("Received message, start consuming", logx.Field("task_info", taskInfo))
		channel := chanType.TypeCodeEn[taskInfo.SendChannel]
		msgType := msgType.TypeCodeEn[taskInfo.MsgType]
		err := task.Submit(ctx, fmt.Sprintf("%s.%s", channel, msgType), task.NewTask(taskInfo, t.svcCtx))
		if err != nil {
			logx.WithContext(ctx).Errorw("submit task err",
				logx.Field("content", string(m.Value)),
				logx.Field("err", err))
			return err
		}
	}
	return nil
}

func (t *TaskExecutor) Stop() {
	if t.isGroup {
		t.ConsumerGroup.Stop()
	} else if t.Consumer != nil {
		t.Consumer.Stop()
	}
}
