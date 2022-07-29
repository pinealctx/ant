package gc

import (
	"context"
	"github.com/pinealctx/ant/example/pb"
	"github.com/pinealctx/ant/pkg/fp"
	"github.com/pinealctx/neptune/ulog"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	EchoCommand = &cli.Command{
		Name:   "ge",
		Usage:  "example client grpc echo",
		Flags:  grFlags,
		Action: runEcho,
	}

	ProcCommand = &cli.Command{
		Name:   "gp",
		Usage:  "example client grpc echo",
		Flags:  grFlags,
		Action: runProc,
	}

	GetCommand = &cli.Command{
		Name:   "gg",
		Usage:  "example client grpc echo",
		Flags:  grFlags,
		Action: runGet,
	}

	TouchCommand = &cli.Command{
		Name:   "gt",
		Usage:  "example client grpc echo",
		Flags:  grFlags,
		Action: runTouch,
	}
)

var (
	grFlags = append(fp.FlowFlags, &cli.StringFlag{
		Name:  "url",
		Usage: "grpc server url",
		Value: "127.0.0.1:5222",
	})
)

func runEcho(c *cli.Context) error {
	var fn = func(ctx context.Context, nc pb.ExampleClient) {
		var _, e = nc.Echo(ctx, shareReq)
		if e != nil {
			ulog.Error("echo.error", zap.Error(e))
		}
	}
	return runHandler(c, fn)
}

func runProc(c *cli.Context) error {
	var fn = func(ctx context.Context, nc pb.ExampleClient) {
		var _, e = nc.Proc(ctx, shareReq)
		if e != nil {
			ulog.Error("proc.error", zap.Error(e))
		}
	}
	return runHandler(c, fn)
}

func runGet(c *cli.Context) error {
	var fn = func(ctx context.Context, nc pb.ExampleClient) {
		var _, e = nc.Get(ctx, emptyMsg)
		if e != nil {
			ulog.Error("get.error", zap.Error(e))
		}
	}
	return runHandler(c, fn)
}

func runTouch(c *cli.Context) error {
	var fn = func(ctx context.Context, nc pb.ExampleClient) {
		var _, e = nc.Touch(ctx, emptyMsg)
		if e != nil {
			ulog.Error("touch.error", zap.Error(e))
		}
	}
	return runHandler(c, fn)
}

func runHandler(c *cli.Context, fn func(context.Context, pb.ExampleClient)) error {
	var nc, err = setCli(c)
	if err != nil {
		return err
	}
	var ffn = func(ctx context.Context) {
		fn(ctx, nc)
	}
	return fp.ParseFromCtx(c).Run(ffn)
}

func setCli(c *cli.Context) (pb.ExampleClient, error) {
	var conn, err = grpc.Dial(c.String("url"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return pb.NewExampleClient(conn), nil
}
