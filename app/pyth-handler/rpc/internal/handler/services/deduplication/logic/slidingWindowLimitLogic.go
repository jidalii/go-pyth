package logic

import (
	"context"
	_ "embed"

	"github.com/zeromicro/go-zero/core/logx"
	
	// "github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/limit"

	"pyth-go/app/pyth-common/types"
	dtypes "pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/types"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
	"pyth-go/common/timex"
)
const slidingWindowLimitTag = "SW_"

// SlidingWindowLimitLogic rate-limits the message freq in a given time window
// 
// Currently used to rate limit daily message consumption
type SlidingWindowLimitLogic struct {
	svcCtx *svc.ServiceContext
}

func NewSlidingWindowLimitLogic(svcCtx *svc.ServiceContext) dtypes.LimitLogic {
	return &SlidingWindowLimitLogic{svcCtx: svcCtx}
}

// Return the receiver that reached the limit
// 
// service: deduplication service
// 
// taskInfo: task information
// 
// cfg: deduplication configuration
func (s *SlidingWindowLimitLogic) Filter(
	ctx context.Context,
	service dtypes.IDeduplicationService,
	taskInfo *types.TaskInfo,
	cfg dtypes.DeduplicationCfg) (filterReceiver []string, err error) {
	filterReceiver = make([]string, 0, len(taskInfo.Receiver))
	period := timex.GetDisTodayEnd()
	periodLimiter := limit.NewPeriodLimit(int(period), cfg.Num, s.svcCtx.RedisClient, slidingWindowLimitTag)

	for _, receiver := range taskInfo.Receiver {
		key := service.GetDeduplicationSingleKey(taskInfo, receiver)
		code, err := periodLimiter.TakeCtx(ctx, key)
		if err != nil {
			logx.Errorw("slidingWindowLimitLogic Take ", logx.Field("err", err))
			continue
		}
		// if reached limit, add to filterReceiver
		if code == limit.OverQuota {
			filterReceiver = append(filterReceiver, receiver)
		}
	}
	return
}
