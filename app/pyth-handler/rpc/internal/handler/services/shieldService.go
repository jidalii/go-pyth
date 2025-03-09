package services

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"pyth-go/app/pyth-common/types"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

const (
	NighttNoShield               = 10
	NightShield                  = 20
	NightShieldButNextDaySend    = 30
	NightShieldButNextDaySendKey = "night_shield_send"
	SeondsOfADay                 = 86400
)

type shieldService struct {
	svcCtx *svc.ServiceContext
}

func NewShieldService(svcCtx *svc.ServiceContext) *shieldService {
	return &shieldService{svcCtx: svcCtx}
}

func (s shieldService) Shield(ctx context.Context, taskInfo *types.TaskInfo, now time.Time) {
	if taskInfo.ShieldType == NighttNoShield {
		return
	}
	if isNight(now) {
		if taskInfo.ShieldType == NightShield {
			taskInfo.Receiver = []string{}
			return
		} else if taskInfo.ShieldType == NightShieldButNextDaySend {
			taskInfoBytes, _ := jsonx.Marshal(taskInfo)

			err := s.svcCtx.RedisClient.PipelinedCtx(ctx, func(pipe redis.Pipeliner) error {
				if _, err := s.svcCtx.RedisClient.Lpush(NightShieldButNextDaySendKey, taskInfoBytes); err != nil {
					return err
				}

				if err := s.svcCtx.RedisClient.Expire(NightShieldButNextDaySendKey, SeondsOfADay); err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				logx.Errorw("shieldService: failed to lpush TaskInfo", logx.Field("taskInfo", *taskInfo), logx.Field("err", err))
			}
		} else {
			logx.Errorw("shieldService: received invalid ShieldType", logx.Field("taskInfo", *taskInfo))
		}
	}
}

// Night range: 22:00 - 7:59
func isNight(now time.Time) bool {
	return now.Hour() < 8 || now.Hour() > 21
}
