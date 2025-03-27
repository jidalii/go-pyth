package handlers

import (
	"context"

	
	
	"pyth-go/app/pyth-common/types"
)

type smsHandler struct {
	BaseHandler
}

func NewSmsHandler() IHandler {
	return &smsHandler{}
}

func (h *smsHandler) Do(ctx context.Context, taskInfo types.TaskInfo) error {
	return nil
}

func (h *smsHandler) Recall(ctx context.Context, taskInfo types.TaskInfo) {
}


