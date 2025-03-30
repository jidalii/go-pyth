package main

import (
	"flag"
	"fmt"

	"pyth-go/app/pyth-handler/rpc/internal/config"
	"pyth-go/app/pyth-handler/rpc/internal/handler"
	"pyth-go/app/pyth-handler/rpc/internal/server"
	"pyth-go/app/pyth-handler/rpc/internal/svc"
	"pyth-go/app/pyth-handler/rpc/pb/pythhanlder"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/pythhandler_test.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pythhanlder.RegisterPythHandlerServer(grpcServer, server.NewPythHandlerServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	serviceGroup.Add(s)
	serviceGroup.Add(handler.NewTaskExecutor(ctx, true))


	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	serviceGroup.Start()
}
