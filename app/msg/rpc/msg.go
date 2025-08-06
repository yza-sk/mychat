package main

import (
	"awesomeProject/common/interceptor"
	"awesomeProject/proto/msg"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"

	"awesomeProject/app/msg/rpc/internal/config"
	"awesomeProject/app/msg/rpc/internal/server"
	"awesomeProject/app/msg/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/msg.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		msg.RegisterMessageClientServer(grpcServer, server.NewMessageClientServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 添加拦截器
	s.AddUnaryInterceptors(interceptor.LoggerInterceptor)
	// 禁用显示cpu
	logx.DisableStat()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
