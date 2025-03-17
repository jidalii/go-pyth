package service

import (
	"context"

    "github.com/zeromicro/go-zero/core/jsonx"
    "github.com/spf13/cast"

	"pyth-go/app/pyth-common/types"
    "pyth-go/common/encrypt"
	"pyth-go/common/zutils/arrayUtils"
	dtypes "pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/types"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

// A content deduplication service
//
// Within the N seconds, a user can only receive the same content once
// 
// Use simpleLimitLogic to filter the receivers
type contentDeduplicationService struct {
	svcCtx *svc.ServiceContext
	limit  dtypes.LimitLogic
}

func NewContentDeduplicationService(svcCtx *svc.ServiceContext, limit dtypes.LimitLogic) dtypes.IDeduplicationService {
	return &contentDeduplicationService{
		svcCtx: svcCtx,
		limit:  limit,
	}
}

func (s contentDeduplicationService) DoDeduplication(ctx context.Context, taskInfo *types.TaskInfo, cfg dtypes.DeduplicationCfg) error {
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

func (s contentDeduplicationService) GetDeduplicationSingleKey(taskInfo *types.TaskInfo, receiver string) string {
	content, _ := jsonx.Marshal(taskInfo.ContentModel)
	return encrypt.MD5(cast.ToString(taskInfo.MessageTemplateId) + receiver + string(content))
}
