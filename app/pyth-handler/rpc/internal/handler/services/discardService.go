package services

import (
	"context"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"

	"pyth-go/app/pyth-common/types"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
	"pyth-go/common/zutils/arrayUtils"
)

type discardMessageService struct {
	svcCtx *svc.ServiceContext
}

const discardMessageKey = "discard_msg"

func NewDiscardMessageService(svcCtx *svc.ServiceContext) *discardMessageService {
	return &discardMessageService{svcCtx: svcCtx}
}

func (l discardMessageService) IsDiscard(ctx context.Context, taskInfo *types.TaskInfo) bool {
	discardMessageTemplateIds, err := l.svcCtx.RedisClient.SmembersCtx(ctx, discardMessageKey)
	if err != nil {
		logx.Errorw("discardMessageService smembers ", logx.Field("err", err))
		return false
	}
	if len(discardMessageTemplateIds) == 0 {
		return false
	}
	if arrayUtils.ArrayIn(discardMessageTemplateIds, strconv.FormatInt(taskInfo.MessageTemplateId, 10)) {
		return true
	}
	return false
}
