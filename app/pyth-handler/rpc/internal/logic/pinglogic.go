package logic

import (
	"context"

	"pyth-go/app/pyth-handler/rpc/internal/svc"
	"pyth-go/app/pyth-handler/rpc/pb/pythhanlder"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *pythhanlder.PingRequest) (*pythhanlder.PongResponse, error) {

	return &pythhanlder.PongResponse{}, nil
}
