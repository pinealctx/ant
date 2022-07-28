package main

import (
	"github.com/pinealctx/ant/example/main/cmd/gc"
	"github.com/pinealctx/ant/example/main/cmd/na"
	"github.com/pinealctx/ant/example/main/cmd/raw"
	"github.com/pinealctx/neptune/ulog"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
)

func main() {
	var app = &cli.App{
		Name:  "nats rpc example",
		Usage: "use nats to handle rpc",
		Commands: cli.Commands{
			na.SrvCommand,
			na.EchoCommand,
			na.ProcCommand,
			na.GetCommand,
			na.TouchCommand,

			gc.SrvCommand,
			gc.EchoCommand,
			gc.ProcCommand,
			gc.GetCommand,
			gc.TouchCommand,

			raw.SrvCommand,
			raw.EchoCommand,
		},
	}
	var err = app.Run(os.Args)
	if err != nil {
		ulog.Error("app.run", zap.Error(err))
	}
}
