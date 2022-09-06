package service

import (
	v1 "cliset/api/helloworld/v1"
	"cliset/internal/biz"
	"context"
	"github.com/gin-gonic/gin"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}

const tpl = `
<a href='./build' >build</a> <br>
<a href='./cat'>cat</a>
`

func (s *GreeterService) Index(ctx *gin.Context) {
	//markdown := blackfriday.MarkdownCommon([]byte(body))
	//ctx.Data(http.StatusOK, "text/html; charset=utf-8", markdown)
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(200, tpl)
	return
}
func (s *GreeterService) Build(ctx *gin.Context) {
	err := s.uc.Build(ctx.Request.Context())
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	if err != nil {
		ctx.String(200, err.Error())
		return
	}
	ctx.String(200, tpl+"执行完毕")
	return
}
func (s *GreeterService) Cat(ctx *gin.Context) {
	str, _ := s.uc.Cat(ctx.Request.Context())
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(200, tpl+str)
}
