package server

import (
	v1 "cliset/api/helloworld/v1"
	"cliset/internal/conf"
	"cliset/internal/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func gin2(srv *service.GreeterService) *gin.Engine {
	router := gin.Default()
	router.Use(kgin.Middlewares(recovery.Recovery(), customMiddleware))

	//rootGrp := router.Group("/gin")
	//{
	//	userGrp := rootGrp.Group("/user")
	//	userGrp.GET("/sayhi", us.Get)
	//}
	router.GET("/index", srv.Index)
	router.GET("/build", srv.Build)
	router.GET("/cat", srv.Cat)
	//router.GET("/index", func(ctx *gin.Context) {
	//	//markdown := blackfriday.MarkdownCommon([]byte(body))
	//	//ctx.Data(http.StatusOK, "text/html; charset=utf-8", markdown)
	//	ctx.Header("Content-Type", "text/html; charset=utf-8")
	//	ctx.String(200, tpl)
	//	return
	//	name := ctx.Param("name")
	//	if name == "error" {
	//		kgin.Error(ctx, errors.Unauthorized("auth_error", "no authentication"))
	//	} else {
	//		ctx.JSON(200, map[string]string{"welcome": name})
	//	}
	//})
	return router
}
func customMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		if tr, ok := transport.FromServerContext(ctx); ok {
			fmt.Println("operation:", tr.Operation())
		}
		reply, err = handler(ctx, req)
		return
	}
}

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	v1.RegisterGreeterHTTPServer(srv, greeter)
	srv.HandlePrefix("/", gin2(greeter))
	return srv
}
