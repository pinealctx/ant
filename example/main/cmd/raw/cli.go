package raw

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/pinealctx/ant/example/pb"
	"github.com/pinealctx/neptune/ulog"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	EchoCommand = &cli.Command{
		Name:   "re",
		Usage:  "example client raw echo",
		Flags:  shareNatsFlags,
		Action: runEcho,
	}
)

var (
	sharedReq      = &pb.Hello{Halo: "halo"}
	shareNatsFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  "url",
			Usage: "nats server url",
			Value: "nats://127.0.0.1:4222",
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
	var fn = func(ctx context.Context, nc *nats.Conn) {
		var _, err = nc.RequestWithContext(ctx, _EchoTopic.next(), []byte("halo"))
		if err != nil {
			ulog.Error("echo.error", zap.Error(err))
		}
	}
	return runHandler(c, fn)
}

func runHandler(c *cli.Context, fn func(context.Context, *nats.Conn)) error {
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
				fn(selfCtx, cc.conn)
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
			fn(ctx, cc.conn)
		}
		var tb = time.Now()
		var diff = tb.Sub(ta)
		ulog.Info("sync use time", zap.String("total", fmt.Sprintf("%+v", diff)),
			zap.String("average", fmt.Sprintf("%+v", diff/time.Duration(cc.count))))
	}
	return nil
}

type cliSetT struct {
	conn       *nats.Conn
	count      int
	async      bool
	inheritCtx bool
	timeout    time.Duration
}

func setCli(c *cli.Context) (*cliSetT, error) {
	var conn, err = connectNats(c)
	if err != nil {
		return nil, err
	}
	return &cliSetT{
		conn:       conn,
		count:      c.Int("count"),
		async:      c.Bool("async"),
		inheritCtx: c.Bool("ctx"),
		timeout:    c.Duration("tc"),
	}, nil
}
