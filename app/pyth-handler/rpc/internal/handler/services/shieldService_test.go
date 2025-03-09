package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/jsonx"

	// "github.com/zeromicro/go-zero/core/stores/redis"

	"pyth-go/app/pyth-common/types"
	"pyth-go/app/pyth-handler/rpc/internal/config"
	"pyth-go/app/pyth-handler/rpc/internal/handler/services"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

func TestShieldService(t *testing.T) {
	// load config
	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()

	// init service
	shieldService := services.NewShieldService(svcCtx)

	testCases := []struct {
		name       string
		now        time.Time
		shieldType int
	}{
		{
			name:       "Not night",
			now:        time.Date(2021, 1, 1, 13, 0, 0, 0, time.Local),
			shieldType: services.NightShield,
		},
		{
			name:       "nightShield",
			now:        time.Date(2021, 1, 1, 3, 0, 0, 0, time.Local),
			shieldType: services.NightShield,
		},
		{
			name:       "NighttNoShield",
			now:        time.Date(2021, 1, 1, 3, 0, 0, 0, time.Local),
			shieldType: services.NighttNoShield,
		},
		{
			name:       "NightShieldButNextDaySend",
			now:        time.Date(2021, 1, 1, 3, 0, 0, 0, time.Local),
			shieldType: services.NightShieldButNextDaySend,
		},
		{
			name:       "NightShieldButNextDaySend",
			now:        time.Date(2021, 1, 1, 23, 0, 0, 0, time.Local),
			shieldType: services.NightShieldButNextDaySend,
		},
	}

	for i := 0; i < len(testCases); i++ {
		testCase := testCases[i]
		taskInfo := &types.TaskInfo{
			MessageTemplateId: int64(i),
			ShieldType:        testCase.shieldType,
			ContentModel:      testCase.name,
		}
		taskInfoBytes, _ := jsonx.Marshal(taskInfo)
		taskInfoStr := string(taskInfoBytes)
		shieldService.Shield(ctx, taskInfo, testCase.now)
		resp, _ := svcCtx.RedisClient.RpopCtx(ctx, services.NightShieldButNextDaySendKey)

		if isNight(testCase.now) {
			if taskInfo.ShieldType == services.NightShield {
				assert.Equal(t, "", resp)
			} else if taskInfo.ShieldType == services.NightShieldButNextDaySend {
				assert.Equal(t, taskInfoStr, resp)
				t.Logf("added msg{shieldType %d}: {%s}", services.NightShieldButNextDaySend, resp)
			} else {
				assert.Equal(t, "", resp)
			}
		} else {
			assert.Equal(t, 0, len(resp))
		}
	}
}

// !!! Attention: before running the test above, make sure this `isNight` function is defined in the same file
func isNight(now time.Time) bool {
	return now.Hour() < 8 || now.Hour() > 21
}
