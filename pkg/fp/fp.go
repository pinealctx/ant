package fp

import (
	"context"
	"fmt"
	"github.com/pinealctx/neptune/ulog"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	totalTimeoutKey  = "atc"
	countKey         = "count"
	asyncKey         = "async"
	inheritedCtxKey  = "icx"
	singleTimeoutKey = "stc"
	concurrencyKey   = "cl"
)

var (
	FlowFlags = []cli.Flag{
		&cli.DurationFlag{
			Name:  totalTimeoutKey,
			Usage: "all rpc total timeout",
			Value: time.Minute * 30,
		},
		&cli.BoolFlag{
			Name:  asyncKey,
			Usage: "set async go routine or not",
		},
		&cli.IntFlag{
			Name:  countKey,
			Usage: "call count",
			Value: 100000,
		},
		&cli.BoolFlag{
			Name:  inheritedCtxKey,
			Usage: "inherit context or not, if true, each rpc will use its own context inherited from root",
		},
		&cli.DurationFlag{
			Name:  singleTimeoutKey,
			Usage: "single rpc timeout, if icx set true, each rpc will set this as timeout",
			Value: time.Second * 6,
		},
		&cli.IntFlag{
			Name:  concurrencyKey,
			Usage: "use semaphore to control call go routine, this is the concurrency go routine number",
			Value: 0,
		},
	}
)

// FlowP 一组调用，go routine/semaphore/sequence test
type FlowP struct {
	TotalTimeout   time.Duration
	Count          int
	Async          bool
	InheritedCtx   bool
	Timeout        time.Duration
	ConcurrencyNum int
}

// Run : run multiple rpc
func (p *FlowP) Run(fn func(context.Context)) error {
	p.logInfo()
	var ctx, cancel = context.WithTimeout(context.Background(), p.TotalTimeout)
	defer cancel()

	if !p.Async {
		p.runSync(ctx, fn)
	} else {
		if p.ConcurrencyNum > 0 {
			p.runAsyncWithLimit(ctx, fn)
		} else {
			p.runAsync(ctx, fn)
		}
	}
	return nil
}

// run async
func (p *FlowP) runAsync(ctx context.Context, fn func(context.Context)) {
	var wg sync.WaitGroup
	wg.Add(p.Count)
	var t1 = time.Now()

	var nfn func()
	if p.InheritedCtx {
		nfn = func() {
			defer wg.Done()
			var sCtx, cancel = context.WithTimeout(ctx, p.Timeout)
			defer cancel()
			fn(sCtx)
		}
	} else {
		nfn = func() {
			defer wg.Done()
			fn(ctx)
		}
	}

	for i := 0; i < p.Count; i++ {
		go nfn()
	}
	wg.Wait()
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	var logKey string
	if p.InheritedCtx {
		logKey = "async with same context use time"
	} else {
		logKey = "async with inherited context use time"
	}

	ulog.Info(logKey, zap.String("total", fmt.Sprintf("%+v", dur)),
		zap.String("average", fmt.Sprintf("%+v", dur/time.Duration(p.Count))))
}

// run async with go routine limit
func (p *FlowP) runAsyncWithLimit(ctx context.Context, fn func(ctx context.Context)) {
	var wg sync.WaitGroup
	wg.Add(p.ConcurrencyNum)
	var t1 = time.Now()

	var counter = NewCounter(p.Count)
	var nfn func()
	if p.InheritedCtx {
		nfn = func() {
			var sCtx, cancel = context.WithTimeout(ctx, p.Timeout)
			defer cancel()
			fn(sCtx)
		}
	} else {
		nfn = func() {
			fn(ctx)
		}
	}

	for i := 0; i < p.ConcurrencyNum; i++ {
		go func() {
			defer wg.Done()
			for counter.Acquire() {
				nfn()
			}
		}()
	}

	wg.Wait()
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	var logKey string
	if p.InheritedCtx {
		logKey = "run async limit with same context use time"
	} else {
		logKey = "run async limit with inherited context use time"
	}

	ulog.Info(logKey, zap.String("total", fmt.Sprintf("%+v", dur)),
		zap.String("average", fmt.Sprintf("%+v", dur/time.Duration(p.Count))))
}

// run sync
func (p *FlowP) runSync(ctx context.Context, fn func(context.Context)) {
	var t1 = time.Now()
	for i := 0; i < p.Count; i++ {
		fn(ctx)
	}
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	ulog.Info("sync call time", zap.String("total", fmt.Sprintf("%+v", dur)),
		zap.String("average", fmt.Sprintf("%+v", dur/time.Duration(p.Count))))
}

// log flow info
func (p *FlowP) logInfo() {
	ulog.Info("current-setting",
		zap.String("total-timeout", fmt.Sprintf("%v", p.TotalTimeout)),
		zap.Int("count", p.Count),
		zap.Bool("async", p.Async),
		zap.Bool("use-single-context", p.InheritedCtx),
		zap.String("single-timeout", fmt.Sprintf("%v", p.Timeout)),
		zap.Int("concurrency-limit", p.ConcurrencyNum),
	)
}

// ParseFromCtx parse a flow control from command line context
func ParseFromCtx(c *cli.Context) *FlowP {
	return &FlowP{
		TotalTimeout:   c.Duration(totalTimeoutKey),
		Count:          c.Int(countKey),
		Async:          c.Bool(asyncKey),
		InheritedCtx:   c.Bool(inheritedCtxKey),
		Timeout:        c.Duration(singleTimeoutKey),
		ConcurrencyNum: c.Int(concurrencyKey),
	}
}

type Counter struct {
	count  int
	locker sync.Mutex
}

func NewCounter(count int) *Counter {
	return &Counter{count: count}
}

func (c *Counter) Acquire() bool {
	var ok bool
	c.locker.Lock()
	if c.count > 0 {
		ok = true
	}
	c.count--
	c.locker.Unlock()
	return ok
}
