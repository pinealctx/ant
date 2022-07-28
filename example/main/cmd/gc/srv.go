package gc

import (
	"context"
	"github.com/pinealctx/ant/example/pb"
	"github.com/pinealctx/neptune/ulog"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

var (
	emptyMsg = &emptypb.Empty{}
	shareReq = &pb.Hello{}
)

var SrvCommand = &cli.Command{
	Name:  "gs",
	Usage: "example grpc service",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "url",
			Usage: "grpc server url",
			Value: ":5222",
		},
	},
	Action: runService,
}

func runService(c *cli.Context) error {
	var url = c.String("url")
	var ln, err = net.Listen("tcp", url)
	if err != nil {
		return err
	}
	var gs = grpc.NewServer()
	pb.RegisterExampleServer(gs, &RPC{})
	ulog.Warn("example grpc service running", zap.String("url", url))
	err = gs.Serve(ln)
	ulog.Error("grpc.run.error", zap.Error(err))
	return err
}

type RPC struct {
	pb.UnimplementedExampleServer
}

func (g *RPC) Echo(ctx context.Context, hello *pb.Hello) (*pb.Hello, error) {
	return hello, nil
}

func (g *RPC) Proc(ctx context.Context, hello *pb.Hello) (*emptypb.Empty, error) {
	return emptyMsg, nil
}

func (g *RPC) Get(ctx context.Context, empty *emptypb.Empty) (*pb.Hello, error) {
	return shareReq, nil
}

func (g *RPC) Touch(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return emptyMsg, nil
}
