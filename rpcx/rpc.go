package rpcx

import (
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/liblpc"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
)

type RPC struct {
	ioLoop          *liblpc.IOEvtLoop
	rcpFuncMap      map[string]*rpcFunc
	promiseGroup    *std.PromiseGroup
	lock            *sync.RWMutex
	startFlag       int32
	callableCloseCb func(callable Callable)
	middlewares     []MiddlewareFunc
}

const RpcLoopDefaultBufferSize = 1024 * 1024 * 4

func New() (*RPC, error) {
	loop, err := liblpc.NewIOEvtLoop(RpcLoopDefaultBufferSize)
	if err != nil {
		return nil, err
	}
	return &RPC{
		ioLoop:       loop,
		rcpFuncMap:   make(map[string]*rpcFunc),
		promiseGroup: std.NewPromiseGroup(),
		lock:         &sync.RWMutex{},
		startFlag:    0,
		middlewares:  make([]MiddlewareFunc, 0),
	}, nil
}
func (this *RPC) OnCallableClosed(cb func(callable Callable)) {
	this.callableCloseCb = cb
}

func (this *RPC) Use(m MiddlewareFunc) {
	this.middlewares = append(this.middlewares, m)
}

func (this *RPC) buildCallChain(direct Direction, h HandleFunc) HandleFunc {
	switch direct {
	case In:
		{
			for i := 0; i < len(this.middlewares); i++ {
				h = this.middlewares[i](h)
			}
		}
	case Out:
		{
			for i := len(this.middlewares) - 1; i > 0; i-- {
				h = this.middlewares[i](h)
			}
		}
	default:
		std.Assert(false, "unknown direction")
	}
	return h
}

func (this *RPC) Loop() liblpc.EventLoop {
	return this.ioLoop
}
func (this *RPC) getFunc(name string) *rpcFunc {
	this.lock.RLock()
	defer this.lock.RUnlock()
	fn, ok := this.rcpFuncMap[name]
	if !ok {
		return nil
	}
	return fn
}

func (this *RPC) RegFuncWithName(fname string, f interface{}) {
	fv, ok := f.(reflect.Value)
	if !ok {
		fv = reflect.ValueOf(f)
	}
	std.Assert(fv.Kind() == reflect.Func, "f not func!")
	fvType := fv.Type()
	//check in/out param
	checkInParam(fvType)
	checkOutParam(fvType)
	//
	this.lock.Lock()
	defer this.lock.Unlock()
	//
	this.rcpFuncMap[fname] = &rpcFunc{
		name:      fname,
		fun:       fv,
		inP0Type:  fvType.In(1),
		outP0Type: fvType.Out(0),
	}
}

func (this *RPC) RegFun(f interface{}) {
	fv, ok := f.(reflect.Value)
	if !ok {
		fv = reflect.ValueOf(f)
	}
	std.Assert(fv.Kind() == reflect.Func, "f not func!")
	fname := getFuncName(fv)
	this.RegFuncWithName(fname, fv)
}

func (this *RPC) Start() {
	if atomic.CompareAndSwapInt32(&this.startFlag, 0, 1) {
		go this.ioLoop.Run()
	}
}

func (this *RPC) Close() error {
	this.ioLoop.Break()
	return this.ioLoop.Close()
}

func (this *RPC) newCallable(stream *liblpc.BufferedStream, userData interface{}, m []MiddlewareFunc) *rpcCli {
	s := &rpcCli{
		stream: stream,
		ctx:    this,
		mid:    make([]MiddlewareFunc, 0),
	}
	//
	s.mid = append(s.mid, m...)
	//
	s.SetUserData(userData)
	s.stream.SetUserData(s)
	//
	return s
}

func (this *RPC) NewConnCallable(fd int, userData interface{}, m ...MiddlewareFunc) Callable {
	stream := liblpc.NewBufferedConnStream(this.ioLoop, fd, this.genericRead)
	pCall := this.newCallable(stream, userData, m)
	pCall.start()
	return pCall
}

type ClientCallableOnConnect = func(callable Callable, err error)

func (this *RPC) NewClientCallable(fd int, userData interface{}, m ...MiddlewareFunc) (cancelFn func(), future std.Future) {
	cliStream := liblpc.NewBufferedClientStream(this.ioLoop, fd, this.genericRead)
	pCall := this.newCallable(cliStream, userData, m)
	promise := std.NewPromise()
	cliStream.SetOnConnect(func(sw liblpc.StreamWriter, err error) {
		if err != nil {
			promise.DoneData(err, nil)
		} else {
			promise.DoneData(nil, pCall)
		}
	})
	pCall.start()
	return func() {
		_ = pCall.Close()
	}, promise.GetFuture()
}

const kMaxRpcMsgBodyLen = 1024 * 1024 * 32

func (this *RPC) genericRead(sw liblpc.StreamWriter, buf std.ReadableBuffer, err error) {
	if err != nil {
		log.Println("RPC READ ERROR ", err)
		std.CloseIgnoreErr(sw)
		if this.callableCloseCb != nil {
			callable := sw.GetUserData().(Callable)
			this.callableCloseCb(callable)
		}
		return
	}
	for {
		rawMsg, err := decodeRpcMsg(buf, kMaxRpcMsgBodyLen)
		if err != nil {
			break
		}
		isReq := rawMsg.Type == rpcReqMsg
		if isReq {
			go this.handleReq(sw, rawMsg)
		} else {
			this.handleAck(rawMsg)
		}
	}
}

func (this *RPC) handleAck(inMsg *rpcRawMsg) {
	// log.Println("RECV ACK id -> ", inMsg.Id)
	this.promiseGroup.DonePromise(std.PromiseId(inMsg.Id), inMsg.GetError(), inMsg.Data)
}

var errRpcFuncNotFound = errors.New("rpc func not found")

func (this *RPC) lastWriteFn(outMsg *rpcRawMsg, ctx Context) {
	err := ctx.Error()
	if err != nil {
		outMsg.SetError(err)
	} else {
		outBytes, err := gRpcSerialization.Marshal(ctx.Response())
		if err != nil {
			outMsg.SetError(err)
		} else {
			outMsg.Data = outBytes
		}
	}
}

var gRpcSerialization = std.MsgPackSerialization

func (this *RPC) handleReq(sw liblpc.StreamWriter, inMsg *rpcRawMsg) {
	outMsg := &rpcRawMsg{
		Id:         inMsg.Id,
		MethodName: inMsg.MethodName,
		Type:       rpcAckMsg,
	}
	// log.Println("RECV REQ id -> ", inMsg.Id)
	fn := this.getFunc(inMsg.MethodName)
	if fn != nil {
		cli := sw.GetUserData().(*rpcCli)
		ctx := newContext(cli, inMsg)
		ctx.SetHeaders(inMsg.Headers)
		//
		h := fn.buildInvoke(ctx)
		h = this.buildCallChain(In, h)
		err := h(ctx)
		if err != nil {
			ctx.SetError(err)
		}
		if err == nil && ctx.Direction() == In {
			h = this.buildCallChain(Out, h)
			err = h(ctx)
			if err != nil {
				ctx.SetError(err)
			}
		}
		if ctx.Error() != nil {
			outMsg.SetError(ctx.Error())
		} else {
			err = outMsg.SetData(ctx.Response())
			if err != nil {
				outMsg.SetError(err)
			}
		}
	} else {
		outMsg.SetError(errRpcFuncNotFound)
	}

	sendBytes, err := encodeRpcMsg(outMsg)
	if err != nil {
		log.Printf("RPC handle REQ Id -> %s, error -> %v", inMsg.Id, err)
		return // encode rpcMsg failed
	}
	sw.Write(sendBytes, false)
	//log.Println("RPC ACK REQ Id -> ", inMsg.Id)
}
