package logic

import (
	"context"
	"time"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"pyth-go/app/pyth-common/types"
	dtypes "pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/types"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

const simpleLimitTag = "SP_"

// Simple limit logic is used to rate limit the daily message frequency
//
// Logic: user would be filtered if the number of messages received of a channel in a day exceeds `cfg.Num`
//
// Tech: Redis mGet + pipeline set
//
// CAUTION: Whenever the key is updated, the expiration time would be RESET
type SimpleLimitLogic struct {
	svcCtx *svc.ServiceContext
}

func NewSimpleLimitLogic(svcCtx *svc.ServiceContext) dtypes.LimitLogic {
	return &SimpleLimitLogic{svcCtx: svcCtx}
}

func (s *SimpleLimitLogic) Filter(
	ctx context.Context,
	service dtypes.IDeduplicationService,
	taskInfo *types.TaskInfo,
	cfg dtypes.DeduplicationCfg) (dupReceivers []string, err error) {
	dupReceivers = make([]string, 0, len(taskInfo.Receiver))
	validToPutReceivers := make(map[string]string, len(taskInfo.Receiver))

	keys := getAllDeduplicationKeys(service, taskInfo, simpleLimitTag)
	
	values, err := s.svcCtx.RedisClient.MgetCtx(ctx, keys...)
	if err != nil {
		logx.Errorw("simpleLimitLogic Filter MgetCtx failed ", logx.Field("err", err))
		return dupReceivers, err
	}

	keyValues := make(map[string]string, len(taskInfo.Receiver))
	for i := range keys {
		keyValues[keys[i]] = values[i]
	}

	for i, receiver := range taskInfo.Receiver {
		key, val := keys[i], values[i]
		if cast.ToInt(val) >= cfg.Num {
			dupReceivers = append(dupReceivers, receiver)
		} else {
			validToPutReceivers[receiver] = key
		}
	}

	err = s.mPutInRedis(ctx, validToPutReceivers, keyValues, cfg.Time)
	if err != nil {
		logx.Errorw("simpleLimitLogic mPutInRedis failed ", logx.Field("err", err))
		return dupReceivers, err
	}

	return dupReceivers, nil
}

// update keys' values in redis
// validToPutReceivers: receiver -> key, record the valid receivers only
// keyValues: key -> value, record the current values of ALL keys
func (s SimpleLimitLogic) mPutInRedis(ctx context.Context, validToPutReceivers, keyValues map[string]string, deduplicationTime int64) error {
	updateDailyCounter(&keyValues, validToPutReceivers)

	if len(keyValues) > 0 {
		return s.svcCtx.RedisClient.PipelinedCtx(ctx, func(pipeliner redis.Pipeliner) error {
			err := pipeliner.MSet(ctx, keyValues).Err()
			if err != nil {
				return err
			}
			for key := range keyValues {
				err = pipeliner.Expire(ctx, key, time.Duration(deduplicationTime)*time.Second).Err()
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	return nil
}

func updateDailyCounter(kv *map[string]string, validToPutReceivers map[string]string) {
	for _, value := range validToPutReceivers {
		if val, ok := (*kv)[value]; ok {
			(*kv)[value] = cast.ToString(cast.ToInt(val) + 1)
		} else {
			(*kv)[value] = "1"
		}
	}
}
