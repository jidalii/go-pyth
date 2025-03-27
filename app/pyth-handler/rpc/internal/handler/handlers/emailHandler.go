package handlers

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/limit"
	"gopkg.in/gomail.v2"

	"pyth-go/app/pyth-common/account"
	"pyth-go/app/pyth-common/content_models"
	"pyth-go/app/pyth-common/types"
	"pyth-go/app/pyth-handler/model"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
)

const emailHandlerRatelimitKey = "email_handler_rate_limit"

type emailHandler struct {
	BaseHandler
	svcContext *svc.ServiceContext
	limit      *limit.TokenLimiter
}

func NewEmailHandler(svcCtx *svc.ServiceContext) IHandler {
	return emailHandler{
		svcContext: svcCtx,
		limit:      limit.NewTokenLimiter(3, 10, svcCtx.RedisClient, emailHandlerRatelimitKey),
	}
}

func (h emailHandler) Limit(ctx context.Context, taskInfo types.TaskInfo) bool {
	return h.limit.Allow()
}

func (h emailHandler) Do(ctx context.Context, taskInfo types.TaskInfo) (err error) {
	var content content_model.EmailContentModel
	getContentModel(taskInfo.ContentModel, &content)
	m := gomail.NewMessage()

	var acc account.EmailAccount
	// err = accountUtils.GetAccount(ctx, taskInfo.SendAccount, &acc)
	// h.svcContext.SendAccountModel.FindOne(ctx, int64(taskInfo.SendAccount))
	err = GetAccount(ctx, h.svcContext.SendAccountModel, taskInfo.SendAccount, &acc)
	if err != nil {
		return errors.Wrap(err, "emailHandler get account err")
	}

	m.SetHeader("From", m.FormatAddress(acc.Username, "Official"))

	m.SetHeader("To", taskInfo.Receiver...) // receivers

	m.SetHeader("Subject", content.Title)
	// send html format email
	m.SetBody("text/html", content.Content)

	d := gomail.NewDialer(acc.Host, acc.Port, acc.Username, acc.Password)
	if err := d.DialAndSend(m); err != nil {
		return errors.Wrap(err, "emailHandler DialAndSend err")
	}
	return nil
}

func (h emailHandler) Recall(ctx context.Context, taskInfo types.TaskInfo) {

}

func GetAccount(ctx context.Context, sendAccountModel model.SendAccountModel, accountId int64, v interface{}) error {
	model, err := sendAccountModel.FindOne(ctx, accountId)
	if err != nil {
		return err
	}
	if model == nil {
		return fmt.Errorf("failed to fetch account: accountId-%d", accountId)
	}
	err = jsonx.Unmarshal([]byte(model.Config), &v)
	return err
}
