package pb

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/pinealctx/neptune/errorx"
	"github.com/pinealctx/neptune/mpb"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

const (
	SlotSize = 1024
)

var (
	_EchoTopic  *topicT
	_ProcTopic  *topicT
	_GetTopic   *topicT
	_TouchTopic *topicT
	emptyMsg    = &emptypb.Empty{}
)

type topicT struct {
	tps  []string
	rr   atomic.Uint64
	slot uint64
}

func (x *topicT) next() string {
	var n = x.rr.Inc()
	return x.tps[n%x.slot]
}

func newTopicT(service string, method string, slot uint64) *topicT {
	var x = &topicT{}
	x.tps = make([]string, slot)
	x.slot = slot
	for i := 0; i < int(slot); i++ {
		x.tps[i] = service + "." + strconv.Itoa(i) + "." + method
	}
	return x
}

func init() {
	_EchoTopic = newTopicT("Example", "Echo", SlotSize)
	_ProcTopic = newTopicT("Example", "Proc", SlotSize)
	_GetTopic = newTopicT("Example", "Get", SlotSize)
	_TouchTopic = newTopicT("Example", "Touch", SlotSize)
}

type ExampleNatsClient struct {
	conn *nats.Conn
}

// NewExampleNatsClient new Example client of nats
func NewExampleNatsClient(conn *nats.Conn) *ExampleNatsClient {
	return &ExampleNatsClient{
		conn: conn,
	}
}

// Close : close client, be careful, after closed, client can not be used anymore
// You need new another client if you need to use client again
func (c *ExampleNatsClient) Close() error {
	return c.conn.Drain()
}

func (c *ExampleNatsClient) Echo(ctx context.Context, req *Hello) (*Hello, error) {
	var data, err = mpb.MarshalMsg(req)
	if err != nil {
		return nil, err
	}
	var natsMsg *nats.Msg
	natsMsg, err = c.conn.RequestWithContext(ctx, _EchoTopic.next(), data)
	if err != nil {
		return nil, err
	}
	var rspMsg proto.Message
	rspMsg, err = mpb.UnmarshalRPC(natsMsg.Data)
	if err != nil {
		return nil, err
	}
	var rsp, ok = rspMsg.(*Hello)
	if !ok {
		return nil, errorx.NewfWithStack(
			"unknown response type:%+v", rspMsg.ProtoReflect().Descriptor().FullName())
	}
	return rsp, nil
}

func (c *ExampleNatsClient) Proc(ctx context.Context, req *Hello) error {
	var data, err = mpb.MarshalMsg(req)
	if err != nil {
		return err
	}
	var natsMsg *nats.Msg
	natsMsg, err = c.conn.RequestWithContext(ctx, _ProcTopic.next(), data)
	if err != nil {
		return err
	}
	return mpb.UnmarshalEmptyRPC(natsMsg.Data)
}

func (c *ExampleNatsClient) Get(ctx context.Context) (*Hello, error) {
	var data = mpb.MarshalEmpty()
	var natsMsg, err = c.conn.RequestWithContext(ctx, _GetTopic.next(), data)
	if err != nil {
		return nil, err
	}
	var rspMsg proto.Message
	rspMsg, err = mpb.UnmarshalRPC(natsMsg.Data)
	if err != nil {
		return nil, err
	}
	var rsp, ok = rspMsg.(*Hello)
	if !ok {
		return nil, errorx.NewfWithStack(
			"unknown response type:%+v", rspMsg.ProtoReflect().Descriptor().FullName())
	}
	return rsp, nil
}

func (c *ExampleNatsClient) Touch(ctx context.Context) error {
	var data = mpb.MarshalEmpty()
	var natsMsg, err = c.conn.RequestWithContext(ctx, _TouchTopic.next(), data)
	if err != nil {
		return err
	}
	return mpb.UnmarshalEmptyRPC(natsMsg.Data)
}

type ServiceHandler interface {
	Echo(req *Hello) (*Hello, error)
	Proc(req *Hello) error
	Get() (*Hello, error)
	Touch() error
}

type ExampleNatService struct {
	conn    *nats.Conn
	handler ServiceHandler
	cch     chan struct{}
}

func NewExampleNatService(conn *nats.Conn, handler ServiceHandler) *ExampleNatService {
	var ns = &ExampleNatService{
		conn:    conn,
		handler: handler,
		cch:     make(chan struct{}),
	}
	ns.conn.SetClosedHandler(func(*nats.Conn) {
		close(ns.cch)
	})
	return ns
}

func (s *ExampleNatService) Serve() error {
	var err error
	for i := 0; i < SlotSize; i++ {
		_, err = s.conn.Subscribe(_EchoTopic.tps[i], s.echoSub)
		if err != nil {
			return err
		}
		_, err = s.conn.Subscribe(_ProcTopic.tps[i], s.procSub)
		if err != nil {
			return err
		}
		_, err = s.conn.Subscribe(_GetTopic.tps[i], s.getSub)
		if err != nil {
			return err
		}
		_, err = s.conn.Subscribe(_TouchTopic.tps[i], s.touchSub)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ExampleNatService) CloseCh() <-chan struct{} {
	//wait close
	return s.cch
}

func (s *ExampleNatService) Close() error {
	var err = s.conn.Drain()
	if err != nil {
		return err
	}
	return nil
}

func (s *ExampleNatService) echoHandler(msg proto.Message) (proto.Message, error) {
	var req, ok = msg.(*Hello)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"unknown.msg.req:%+v", msg.ProtoReflect().Descriptor().FullName())
	}
	return s.handler.Echo(req)
}

func (s *ExampleNatService) procHandler(msg proto.Message) (proto.Message, error) {
	var req, ok = msg.(*Hello)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"unknown.msg.req:%+v", msg.ProtoReflect().Descriptor().FullName())
	}
	var err = s.handler.Proc(req)
	if err != nil {
		return nil, err
	}
	return emptyMsg, nil
}

func (s *ExampleNatService) getHandler(msg proto.Message) (proto.Message, error) {
	var _, ok = msg.(*emptypb.Empty)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"unknown.msg.req:%+v", msg.ProtoReflect().Descriptor().FullName())
	}
	return s.handler.Get()
}

func (s *ExampleNatService) touchHandler(msg proto.Message) (proto.Message, error) {
	var _, ok = msg.(*emptypb.Empty)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"unknown.msg.req:%+v", msg.ProtoReflect().Descriptor().FullName())
	}
	var err = s.handler.Touch()
	if err != nil {
		return nil, err
	}
	return emptyMsg, nil
}

func (s *ExampleNatService) echoSub(msg *nats.Msg) {
	procProtoMsg(msg, s.echoHandler)
}

func (s *ExampleNatService) procSub(msg *nats.Msg) {
	procProtoMsg(msg, s.procHandler)
}

func (s *ExampleNatService) getSub(msg *nats.Msg) {
	procProtoMsg(msg, s.getHandler)
}

func (s *ExampleNatService) touchSub(msg *nats.Msg) {
	procProtoMsg(msg, s.touchHandler)
}

func procProtoMsg(msg *nats.Msg, fn func(proto.Message) (proto.Message, error)) {
	var pMsg, err = mpb.UnmarshalMsg(msg.Data)
	if err != nil {
		replyErrMsg(msg, status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	var rspMsg proto.Message
	rspMsg, err = fn(pMsg)
	if err != nil {
		replyErrMsg(msg, err)
		return
	}
	replyRspMsg(msg, rspMsg)
}

func replyErrMsg(msg *nats.Msg, err error) {
	var errRsp, e = mpb.MarshalError(err)
	if e != nil {
		ulog.Error("marshal.error.at.response.error.msg", zap.Error(e))
		return
	}
	e = msg.Respond(errRsp)
	if e != nil {
		ulog.Error("reply.error.msg", zap.Error(e))
	}
}

func replyRspMsg(msg *nats.Msg, rspMsg proto.Message) {
	var rsp, err = mpb.MarshalMsg(rspMsg)
	if err != nil {
		ulog.Error("marshal.error.at.response.msg", zap.Error(err))
		return
	}
	err = msg.Respond(rsp)
	if err != nil {
		ulog.Error("reply.msg", zap.Error(err))
	}
}
