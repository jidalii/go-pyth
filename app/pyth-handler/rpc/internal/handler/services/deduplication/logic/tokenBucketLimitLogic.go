package logic

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"pyth-go/app/pyth-common/types"
	dtypes "pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/types"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

const tokenBucketLimitTag = "TB_"

var (
	//go:embed slidingWindowScript.lua
	tokenBucketLuaScript string
	tokenBucketScript    = redis.NewScript(tokenBucketLuaScript)
)

// TokenBucketLimitLogic rates limit the message frequency in a dynamic sliding window
//
// Currently used to deduplicate the identical message to the same receiver within a sliding window time
type TokenBucketLimitLogic struct {
	svcCtx *svc.ServiceContext
}

func NewTokenBucketLimitLogic(svcCtx *svc.ServiceContext) dtypes.LimitLogic {
	return &TokenBucketLimitLogic{svcCtx: svcCtx}
}

// Return the receiver that reached the limit
func (s *TokenBucketLimitLogic) Filter(
	ctx context.Context,
	service dtypes.IDeduplicationService,
	taskInfo *types.TaskInfo,
	cfg dtypes.DeduplicationCfg) ([]string, error) {
	dupReceivers := make([]string, 0)

	for _, receiver := range taskInfo.Receiver {
		key := fmt.Sprintf("%s%s", tokenBucketLimitTag, service.GetDeduplicationSingleKey(taskInfo, receiver))
		resp, err := s.svcCtx.RedisClient.ScriptRunCtx(
			ctx,
			tokenBucketScript,
			[]string{key},
			[]string{
				strconv.FormatInt(cfg.Time, 10),
				strconv.FormatInt(time.Now().Unix(), 10),
				strconv.Itoa(cfg.Num),
				uuid.New().String(),
			},
		)
		if err != nil {
			logx.Errorw("tokenBucketLimitLogic Take ", logx.Field("err", err))
			return dupReceivers, err
		}
		code, ok := resp.(int64)
		if !ok {
			logx.Errorw("tokenBucketLimitLogic Take ", logx.Field("err", "unknown status code"))
			return dupReceivers, fmt.Errorf("unknown status code: %v", resp)
		}
		if code == 1 {
			dupReceivers = append(dupReceivers, receiver)
		}
	}
	return dupReceivers, nil
}
