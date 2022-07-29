package na

import (
	"context"
	"github.com/pinealctx/ant/example/pb"
	"github.com/pinealctx/ant/pkg/fp"
	"github.com/pinealctx/neptune/ulog"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	EchoCommand = &cli.Command{
		Name:   "ne",
		Usage:  "example client echo",
		Flags:  shareNatsFlags,
		Action: runEcho,
	}

	ProcCommand = &cli.Command{
		Name:   "np",
		Usage:  "example client echo",
		Flags:  shareNatsFlags,
		Action: runProc,
	}

	GetCommand = &cli.Command{
		Name:   "ng",
		Usage:  "example client echo",
		Flags:  shareNatsFlags,
		Action: runGet,
	}

	TouchCommand = &cli.Command{
		Name:   "nt",
		Usage:  "example client echo",
		Flags:  shareNatsFlags,
		Action: runTouch,
	}
)

var (
	sharedReq      = &pb.Hello{Halo: "halo"}
	shareNatsFlags = append(fp.FlowFlags, &cli.StringFlag{
		Name:  "url",
		Usage: "nats server url",
		Value: "nats://127.0.0.1:4222",
	})
)

func runEcho(c *cli.Context) error {
	var fn = func(ctx context.Context, nc *pb.ExampleNatsClient) {
		var _, e = nc.Echo(ctx, sharedReq)
		if e != nil {
			ulog.Error("echo.error", zap.Error(e))
		}
	}
	return runHandler(c, fn)
}

func runProc(c *cli.Context) error {
	var fn = func(ctx context.Context, nc *pb.ExampleNatsClient) {
		var e = nc.Proc(ctx, sharedReq)
		if e != nil {
			ulog.Error("proc.error", zap.Error(e))
		}
	}
	return runHandler(c, fn)
}

func runGet(c *cli.Context) error {
	var fn = func(ctx context.Context, nc *pb.ExampleNatsClient) {
		var _, e = nc.Get(ctx)
		if e != nil {
			ulog.Error("get.error", zap.Error(e))
		}
	}
	return runHandler(c, fn)
}

func runTouch(c *cli.Context) error {
	var fn = func(ctx context.Context, nc *pb.ExampleNatsClient) {
		var e = nc.Touch(ctx)
		if e != nil {
			ulog.Error("touch.error", zap.Error(e))
		}
	}
	return runHandler(c, fn)
}

func runHandler(c *cli.Context, fn func(context.Context, *pb.ExampleNatsClient)) error {
	var nc, err = setCli(c)
	if err != nil {
		return err
	}
	var ffn = func(ctx context.Context) {
		fn(ctx, nc)
	}
	return fp.ParseFromCtx(c).Run(ffn)
}

func setCli(c *cli.Context) (*pb.ExampleNatsClient, error) {
	var conn, err = connectNats(c)
	if err != nil {
		return nil, err
	}
	return pb.NewExampleNatsClient(conn), nil
}
