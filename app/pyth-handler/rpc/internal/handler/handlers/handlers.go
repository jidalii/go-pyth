package handlers

import (
	"context"
	"sync"


	"pyth-go/app/pyth-common/enums/chanType"
	"pyth-go/app/pyth-common/types"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

var (
	once          sync.Once
	handlerHolder map[int]IHandler
)

// SetUp 初始化所有handler
func SetUp(svcCtx *svc.ServiceContext) {
	once.Do(func() {
		handlerHolder = map[int]IHandler{
			chanType.Sms:                NewSmsHandler(),
			chanType.Email:              NewEmailHandler(svcCtx),
			// chanType.OfficialAccounts:   NewOfficialAccountHandler(),
			// chanType.EnterpriseWeChat:   NewEnterpriseWeChatHandler(),
			// chanType.DingDing:           NewDingDingRobotHandler(),
			// chanType.DingDingWorkNotice: NewDingDingWorkNoticeHandler(),
		}
	})

}

func GetHandler(sendChannel int) IHandler {
	return handlerHolder[sendChannel]
}

type BaseHandler struct {
}

func (b BaseHandler) Limit(ctx context.Context, taskInfo types.TaskInfo) bool {
	return true
}
