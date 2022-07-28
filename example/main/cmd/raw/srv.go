package raw

import (
	"github.com/nats-io/nats.go"
	"github.com/pinealctx/neptune/ulog"
	"github.com/urfave/cli/v2"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const (
	SlotSize = 1024
)

var (
	_EchoTopic *topicT
)

var SrvCommand = &cli.Command{
	Name:  "rs",
	Usage: "example nats raw service",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "url",
			Usage: "nats server url",
			Value: "nats://127.0.0.1:4222",
		},
	},
	Action: runService,
}

type topicT struct {
	tps  []string
	rr   atomic.Uint64
	slot uint64
}

func (x *topicT) next() string {
	var n = x.rr.Inc()
	return x.tps[n%x.slot]
}

func newTopicT(topic string, slot uint64) *topicT {
	var x = &topicT{}
	x.tps = make([]string, slot)
	x.slot = slot
	for i := 0; i < int(slot); i++ {
		x.tps[i] = topic + "." + strconv.Itoa(i)
	}
	return x
}

func init() {
	_EchoTopic = newTopicT("Example", SlotSize)
}

type NatRawEchoService struct {
	conn *nats.Conn
	cch  chan struct{}
}

func NewNatRawEchoService(conn *nats.Conn) *NatRawEchoService {
	var ns = &NatRawEchoService{
		conn: conn,
		cch:  make(chan struct{}),
	}
	ns.conn.SetClosedHandler(func(*nats.Conn) {
		close(ns.cch)
	})
	return ns
}

func (s *NatRawEchoService) Serve() error {
	var err error
	for i := 0; i < SlotSize; i++ {
		_, err = s.conn.Subscribe(_EchoTopic.tps[i], s.echoSub)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *NatRawEchoService) echoSub(msg *nats.Msg) {
	var err = msg.Respond(msg.Data)
	if err != nil {
		ulog.Error("response error", zap.Error(err))
	}
}

func (s *NatRawEchoService) CloseCh() <-chan struct{} {
	//wait close
	return s.cch
}

func (s *NatRawEchoService) Close() error {
	var err = s.conn.Drain()
	if err != nil {
		return err
	}
	return nil
}

func runService(c *cli.Context) error {
	var conn, err = connectNats(c)
	if err != nil {
		return err
	}
	var service = NewNatRawEchoService(conn)

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
