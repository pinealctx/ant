// Code generated by frgen , DO NOT EDIT.
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
	// SlotSize : slot size is set 1024 in most case.
	// slot is a parallel handler number
	SlotSize = 1024
)

// ExampleNatsClient wrap a rpc client through nats
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

// Touch rpc call
func (c *ExampleNatsClient) Touch(ctx context.Context) error {
	var data = mpb.MarshalEmpty()
	var natsMsg, err = c.conn.RequestWithContext(ctx, _TouchTopic.next(), data)
	if err != nil {
		return err
	}
	return mpb.UnmarshalEmptyRPC(natsMsg.Data)
}

// ServiceHandler nats rpc service interface
// please implement this interface for logic
type ServiceHandler interface {
	// Touch rpc service
	Touch() error
}

// ExampleNatService : actually service is use subscriber function of nats to serve rpc call.
// You need to pass the nats connection and implement an interface of "ServiceHandler" then pass in.
type ExampleNatService struct {
	conn    *nats.Conn
	handler ServiceHandler
	cch     chan struct{}
}

// NewExampleNatService : pass nats connection and an implemented logic of interface "ServiceHandler"
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

// WaitClose after server run, each subscriber has its own go routine, the channel can be use to notify all
// go routine exist. Actually, you can simply use WaitClose, but if there are multiple service, you can use
// select to handle multiple exist
func (s *ExampleNatService) CloseCh() <-chan struct{} {
	//wait close
	return s.cch
}

// Close : close the nats connection can let all subscriber exist
func (s *ExampleNatService) Close() error {
	var err = s.conn.Drain()
	if err != nil {
		return err
	}
	return nil
}

// Serve : run the service, actually, it shall new multiple go routine, please use
// Close/WaitClose or CloseCh to control the service lifecycle
func (s *ExampleNatService) Serve() error {
	var err error
	for i := 0; i < SlotSize; i++ {
		_, err = s.conn.Subscribe(_TouchTopic.tps[i], func(msg *nats.Msg) {
			procProtoMsg(msg, s._TouchHandler)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// WaitClose after server run, each subscriber has its own go routine, use this function to wait for all
// go routine exist
func (s *ExampleNatService) WaitClose() {
	<-s.cch
}

// internal wrapper function of Touch
func (s *ExampleNatService) _TouchHandler(msg proto.Message) (proto.Message, error) {
	var _, ok = msg.(*emptypb.Empty)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument,
			"unknown.msg.req:%+v", msg.ProtoReflect().Descriptor().FullName())
	}
	var err = s.handler.Touch()
	if err != nil {
		return nil, err
	}
	return _EmptyMsg, nil
}

var (
	// topics here actually are rpc tunnels
	_TouchTopic *topicT
	_EmptyMsg   = &emptypb.Empty{}
)

// topic wrapper, which can iterate rpc as RR load balance
type topicT struct {
	tps  []string
	rr   atomic.Uint64
	slot uint64
}

// next topic to handle current rpc
func (x *topicT) next() string {
	var n = x.rr.Inc()
	return x.tps[n%x.slot]
}

// set up a topic wrapper
func newTopicT(service string, method string, slot uint64) *topicT {
	var x = &topicT{}
	x.tps = make([]string, slot)
	x.slot = slot
	for i := 0; i < int(slot); i++ {
		x.tps[i] = service + "." + strconv.Itoa(i) + "." + method
	}
	return x
}

// init function to setup topics as rpc tunnels
func init() {
	// setup topic wrappers
	_TouchTopic = newTopicT("Example", "Touch", SlotSize)
}

// handle nats message as rpc call(proto.Message in, proto.Message out).
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

// marshal err and reply back
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

// marshal response message back
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
