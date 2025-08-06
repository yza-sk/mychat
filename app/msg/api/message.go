package main

import (
	"awesomeProject/app/msg/api/internal/config"
	"awesomeProject/app/msg/api/internal/handler"
	"awesomeProject/app/msg/api/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/message.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logx.DisableStat() // 禁用显示cpu
	//// http升级为websocket
	//server.AddRoute(
	//	rest.Route{
	//		Method: http.MethodGet,
	//		Path:   "/ws",
	//		Handler: func(w http.ResponseWriter, r *http.Request) {
	//			err := logic.ServeWs(ctx, w, r)
	//			xresp.Response(r, w, nil, err)
	//		},
	//	},
	//	rest.WithJwt(ctx.Config.JwtAuth.AccessSecret),
	//)
	//// 开启协程, 专门从MQ中获取消息, 发给对应的群
	//go logic.ConsumeMsgFromMQ(ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
