package na

import (
	"github.com/nats-io/nats.go"
	"github.com/pinealctx/ant/example/pb"
	"github.com/pinealctx/neptune/ulog"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var SrvCommand = &cli.Command{
	Name:  "ns",
	Usage: "example nats service",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "url",
			Usage: "nats server url",
			Value: "nats://127.0.0.1:4222",
		},
	},
	Action: runService,
}

func runService(c *cli.Context) error {
	var conn, err = connectNats(c)
	if err != nil {
		return err
	}
	var service = pb.NewExampleNatService(conn, &RPC{})

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	err = service.Serve()
	if err != nil {
		return err
	}
	var ch = service.CloseCh()
	ulog.Warn("example nats service running")
	select {
	case <-ch:
		ulog.Error("example service consumer closed")
	case sg := <-sigChan:
		err = service.Close()
		if err != nil {
			ulog.Error("example service consumer stop with signal",
				zap.Reflect("signal", sg), zap.Error(err))
		} else {
			ulog.Warn("example service consumer stop with signal", zap.Reflect("signal", sg))
		}
		<-ch
		ulog.Warn("wait service consumer close")
	}
	return nil
}

func connectNats(c *cli.Context) (*nats.Conn, error) {
	var url = c.String("url")
	ulog.Info("nats", zap.String("url", url))
	var conn, err = nats.Connect(url, nats.NoEcho(), nats.MaxReconnects(-1))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type RPC struct {
}

func (R RPC) Echo(req *pb.Hello) (*pb.Hello, error) {
	return req, nil
}

func (R RPC) Proc(req *pb.Hello) error {
	return nil
}

func (R RPC) Get() (*pb.Hello, error) {
	return &pb.Hello{
		Halo: "get.rsp",
	}, nil
}

func (R RPC) Touch() error {
	return nil
}
