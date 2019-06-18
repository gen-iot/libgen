package rpcx

import (
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/liblpc"
	"reflect"
	"sync"
	"sync/atomic"
)

type RPC struct {
	ioLoop       *liblpc.IOEvtLoop
	rcpFuncMap   map[string]*rpcFunc
	promiseGroup *std.PromiseGroup
	lock         *sync.RWMutex
	startFlag    int32
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
	}, nil
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

func (this *RPC) RegFun(f interface{}) {
	fv, ok := f.(reflect.Value)
	if !ok {
		fv = reflect.ValueOf(f)
	}
	std.Assert(fv.Kind() == reflect.Func, "f not func!")
	fname := getFuncName(fv)
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
		inP0Type:  fvType.In(0),
		outP0Type: fvType.Out(0),
	}
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

func (this *RPC) NewCallable(fd int, userData interface{}) Callable {
	s := &apiClient{
		FdBufferedStream: liblpc.NewFdBufferedStream(this.ioLoop, fd, this.genericRead),
		ctx:              this,
	}
	s.SetUserData(userData)
	s.Start()
	return s
}

const kMaxRpcMsgBodyLen = 1024 * 1024 * 32

func (this *RPC) genericRead(sw liblpc.StreamWriter, buf std.ReadableBuffer, err error) {
	if err != nil {
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
	this.promiseGroup.DonePromise(std.PromiseId(inMsg.Id), inMsg.GetError(), inMsg.Data)
}

func (this *RPC) handleReq(sw liblpc.StreamWriter, inMsg *rpcRawMsg) {
	fn := this.getFunc(inMsg.MethodName)
	if fn == nil {
		return // fn not found
	}
	outBytes, err := fn.Call(inMsg.Data)
	outMsg := &rpcRawMsg{
		Id:         inMsg.Id,
		MethodName: inMsg.MethodName,
		Type:       rpcAckMsg,
	}
	if err != nil {
		outMsg.SetError(err)
	} else {
		outMsg.Data = outBytes
	}
	sendBytes, err := encodeRpcMsg(outMsg)
	if err != nil {
		return // encode rpcMsg failed
	}
	sw.Write(sendBytes, false)
}
