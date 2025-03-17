package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"pyth-go/app/pyth-common/types"
	"pyth-go/app/pyth-handler/model"
	"pyth-go/app/pyth-handler/rpc/internal/config"
	"pyth-go/app/pyth-handler/rpc/internal/handler/services"
	dtypes "pyth-go/app/pyth-handler/rpc/internal/handler/services/deduplication/types"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

var messageTemplateTest1 = &model.MessageTemplate{
	Id:                  1,
	Name:                "test1",
	AuditStatus:         20,
	MsgStatus:           30,
	IdType:              10,
	SendChannel:         20,
	TemplateType:        10,
	MsgType:             10,
	ShieldType:          10,
	MsgContent:          "test1",
	SendAccount:         1,
	Creator:             "test",
	Updator:             "test",
	Auditor:             "test",
	Team:                "test",
	Proposer:            "test",
	IsDeleted:           0,
	Created:             time.Now().Unix(),
	Updated:             time.Now().Unix(),
	DeduplicationConfig: "{\"deduplication_10\":{\"time\":86400,\"num\":2},\"deduplication_20\":{\"time\":30,\"num\":10}}",
	TemplateSn:          "test",
}
var messageTemplateTest2 = &model.MessageTemplate{
	Id:                  2,
	Name:                "test2",
	AuditStatus:         20,
	MsgStatus:           30,
	IdType:              10,
	SendChannel:         20,
	TemplateType:        10,
	MsgType:             10,
	ShieldType:          10,
	MsgContent:          "test2",
	SendAccount:         1,
	Creator:             "test",
	Updator:             "test",
	Auditor:             "test",
	Team:                "test",
	Proposer:            "test",
	IsDeleted:           0,
	Created:             time.Now().Unix(),
	Updated:             time.Now().Unix(),
	DeduplicationConfig: "{\"deduplication_10\":{\"time\":86400,\"num\":100},\"deduplication_20\":{\"time\":30,\"num\":20}}",
	TemplateSn:          "test",
}

var messageTemplates = map[int64]*model.MessageTemplate{
	1: messageTemplateTest1,
	2: messageTemplateTest2,
}

func TestDeduplicationService(t *testing.T) {
	// load config
	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()

	// init service
	dedupService, err := services.NewDeduplicationService(svcCtx)
	assert.Nil(t, err)
	defer dedupService.Close()

	dropAndReinitializeTable(t, svcCtx)
	creteTestMessageTemplates(t, svcCtx)
	flushRedis(t, svcCtx)

	testCases := initTestCases(t)

	for _, testCase := range *testCases {
		totalReceivers := len(testCase.taskInfo.Receiver)
		dedupCfg := testCase.dedupConfig
		var limit int
		if dedupCfg["deduplication_10"].Num < dedupCfg["deduplication_20"].Num {
			limit = int(dedupCfg["deduplication_10"].Num)
		} else {
			limit = int(dedupCfg["deduplication_20"].Num)
		}

		taskInfo := testCase.taskInfo
		for i := 0; i < limit; i++ {
			dedupService.Do(ctx, taskInfo)
			assert.Equal(t, totalReceivers, len(taskInfo.Receiver))
		}
		dedupService.Do(ctx, taskInfo)
		assert.Equal(t, 0, len(taskInfo.Receiver))
	}
}

func dropAndReinitializeTable(t *testing.T, svcCtx *svc.ServiceContext) {
	sqlConn := sqlx.NewMysql(svcCtx.Config.DB.DataSource)
	if _, err := sqlConn.Exec("TRUNCATE TABLE `message_template`"); err != nil {
		t.Fatal(err)
	}
}

func creteTestMessageTemplates(t *testing.T, svcCtx *svc.ServiceContext) {
	ctx := context.Background()
	for _, messageTemplate := range messageTemplates {
		if _, err := svcCtx.MessageTemplateModel.Insert(ctx, messageTemplate); err != nil {
			t.Fatal(err)
		}
	}
}

func flushRedis(t *testing.T, svcCtx *svc.ServiceContext) {
	ctx := context.Background()
	if _, err := svcCtx.RedisClient.EvalCtx(ctx, "return redis.call('flushdb')", []string{}); err != nil {
		t.Fatal(err)
	}
}

func initTestCases(t *testing.T) *[]struct {
	taskInfo    *types.TaskInfo
	dedupConfig map[string]dtypes.DeduplicationCfg
} {
	testCases := []struct {
		taskInfo    *types.TaskInfo
		dedupConfig map[string]dtypes.DeduplicationCfg
	}{
		{
			taskInfo: &types.TaskInfo{
				MessageTemplateId: 1,
				ContentModel:      "test-contentModel-Dedup1",
				SendChannel:       10,
				MsgType:           10,
				Receiver:          []string{"alice1", "bob1", "cindy1", "david1"},
			},
		},
		{
			taskInfo: &types.TaskInfo{
				MessageTemplateId: 2,
				ContentModel:      "test-contentModel-Dedup2",
				SendChannel:       20,
				MsgType:           20,
				Receiver:          []string{"alice2", "bob2", "cindy2", "david2"},
			},
		},
	}
	for i := range len(testCases) {
		curConfig := make(map[string]dtypes.DeduplicationCfg)
		err := jsonx.Unmarshal([]byte(messageTemplates[testCases[i].taskInfo.MessageTemplateId].DeduplicationConfig), &curConfig)
		assert.Nil(t, err)
		testCases[i].dedupConfig = curConfig
	}
	return &testCases
}
