package gc

import (
	"context"
	"fmt"
	"github.com/pinealctx/ant/example/pb"
	"github.com/pinealctx/neptune/ulog"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
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
	grFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  "url",
			Usage: "grpc server url",
			Value: "127.0.0.1:5222",
		},
		&cli.BoolFlag{
			Name:  "async",
			Usage: "set async go routine or not",
		},
		&cli.IntFlag{
			Name:  "count",
			Usage: "call count",
			Value: 100000,
		},
		&cli.BoolFlag{
			Name:  "ctx",
			Usage: "inherit context or not",
		},
		&cli.DurationFlag{
			Name:  "tc",
			Usage: "rpc timeout",
			Value: time.Second * 6,
		},
	}
)

func runEcho(c *cli.Context) error {
	var fn = func(ctx context.Context, nc pb.ExampleClient) {
		var _, err = nc.Echo(ctx, shareReq)
		if err != nil {
			ulog.Error("echo.error", zap.Error(err))
		}
	}
	return runHandler(c, fn)
}

func runProc(c *cli.Context) error {
	var fn = func(ctx context.Context, nc pb.ExampleClient) {
		var _, err = nc.Proc(ctx, shareReq)
		if err != nil {
			ulog.Error("proc.error", zap.Error(err))
		}
	}
	return runHandler(c, fn)
}

func runGet(c *cli.Context) error {
	var fn = func(ctx context.Context, nc pb.ExampleClient) {
		var _, err = nc.Get(ctx, emptyMsg)
		if err != nil {
			ulog.Error("get.error", zap.Error(err))
		}
	}
	return runHandler(c, fn)
}

func runTouch(c *cli.Context) error {
	var fn = func(ctx context.Context, nc pb.ExampleClient) {
		var _, err = nc.Touch(ctx, emptyMsg)
		if err != nil {
			ulog.Error("touch.error", zap.Error(err))
		}
	}
	return runHandler(c, fn)
}

func runHandler(c *cli.Context, fn func(context.Context, pb.ExampleClient)) error {
	var cc, err = setCli(c)
	if err != nil {
		return err
	}

	ulog.Info("timeout setting", zap.String("tc", fmt.Sprintf("%+v", cc.timeout)))

	var ctx, cancel = context.WithTimeout(context.Background(), time.Minute*30)
	defer cancel()

	if cc.async {
		var wg sync.WaitGroup
		wg.Add(cc.count)
		var ta = time.Now()
		for i := 0; i < cc.count; i++ {
			go func() {
				defer wg.Done()
				var selfCtx context.Context
				if cc.inheritCtx {
					var selfCancel context.CancelFunc
					selfCtx, selfCancel = context.WithTimeout(ctx, cc.timeout)
					defer selfCancel()
				} else {
					selfCtx = ctx
				}
				fn(selfCtx, cc.nc)
			}()
		}
		wg.Wait()
		var tb = time.Now()
		var diff = tb.Sub(ta)
		ulog.Info("async use time", zap.String("total", fmt.Sprintf("%+v", diff)),
			zap.String("average", fmt.Sprintf("%+v", diff/time.Duration(cc.count))))
	} else {
		var ta = time.Now()
		for i := 0; i < cc.count; i++ {
			fn(ctx, cc.nc)
		}
		var tb = time.Now()
		var diff = tb.Sub(ta)
		ulog.Info("sync use time", zap.String("total", fmt.Sprintf("%+v", diff)),
			zap.String("average", fmt.Sprintf("%+v", diff/time.Duration(cc.count))))
	}
	return nil
}

type cliSetT struct {
	nc         pb.ExampleClient
	count      int
	async      bool
	inheritCtx bool
	timeout    time.Duration
}

func setCli(c *cli.Context) (*cliSetT, error) {
	var conn, err = grpc.Dial(c.String("url"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	var nc = pb.NewExampleClient(conn)
	return &cliSetT{
		nc:         nc,
		count:      c.Int("count"),
		async:      c.Bool("async"),
		inheritCtx: c.Bool("ctx"),
		timeout:    c.Duration("tc"),
	}, nil
}
