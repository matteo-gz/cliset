package biz

import (
	"bytes"
	v1 "cliset/api/helloworld/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"os/exec"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.
type Greeter struct {
	Hello string
}

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	Save(context.Context, *Greeter) (*Greeter, error)
	Update(context.Context, *Greeter) (*Greeter, error)
	FindByID(context.Context, int64) (*Greeter, error)
	ListByHello(context.Context, string) ([]*Greeter, error)
	ListAll(context.Context) ([]*Greeter, error)
	GetDir() string
}

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo GreeterRepo
	log  *log.Helper
	dir  string
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, log: log.NewHelper(logger)}
}

var (
	CmdOut bytes.Buffer
	CmdErr bytes.Buffer
)

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *GreeterUsecase) CreateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}
func (uc *GreeterUsecase) Build(ctx context.Context) error {
	dir := uc.repo.GetDir()
	c := exec.Command("sh")
	in := bytes.NewBuffer(nil)
	c.Stdin = in
	s := "cd " + dir + "  && git clean -d -f && git reset --hard && git pull"
	in.WriteString(s)
	CmdErr = bytes.Buffer{}
	CmdOut = bytes.Buffer{}
	c.Stdout = &CmdOut
	c.Stderr = &CmdErr
	err := c.Start()
	return err
}
func (uc *GreeterUsecase) Cat(ctx context.Context) (string, error) {
	return CmdOut.String() + CmdErr.String(), nil
}
