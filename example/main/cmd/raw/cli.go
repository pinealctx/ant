package raw

import (
	"context"
	"github.com/pinealctx/ant/pkg/fp"
	"github.com/pinealctx/neptune/ulog"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
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
	shareNatsFlags = append(fp.FlowFlags, &cli.StringFlag{
		Name:  "url",
		Usage: "nats server url",
		Value: "nats://127.0.0.1:4222",
	})
)

func runEcho(c *cli.Context) error {
	var nc, err = connectNats(c)
	if err != nil {
		return err
	}
	var fn = func(ctx context.Context) {
		var _, e = nc.RequestWithContext(ctx, _EchoTopic.next(), []byte("halo"))
		if e != nil {
			ulog.Error("echo.error", zap.Error(e))
		}
	}
	return fp.ParseFromCtx(c).Run(fn)
}
