package main

import (
	"flag"
	"fmt"

	"github.com/colinrs/shopjoy/admin/internal/config"
	"github.com/colinrs/shopjoy/admin/internal/handler"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/pkg/response"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/admin-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)
	httpx.SetErrorHandlerCtx(response.ErrHandle)
	httpx.SetOkHandler(response.OKHandle)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
