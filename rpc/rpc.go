package rpc

import (
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/liblpc"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type rpcFunc struct {
	name      string
	fun       reflect.Value
	inP0Type  reflect.Type
	outP0Type reflect.Type
}

func (this *rpcFunc) decodeInParam(data []byte) (interface{}, error) {
	elementType := this.inP0Type
	isPtr := false
	if this.inP0Type.Kind() == reflect.Ptr {
		elementType = this.inP0Type.Elem()
		isPtr = true
	}
	newOut := reflect.New(elementType).Interface()
	err := std.MsgpackUnmarshal(data, newOut)
	if err != nil {
		return nil, err
	}
	if !isPtr {
		newOut = reflect.ValueOf(newOut).Elem().Interface()
	}
	return newOut, nil
}

func (this *rpcFunc) Call(inBytes []byte) (outBytes []byte, err error) {
	defer func() {
		panicErr := recover()
		if panicErr == nil {
			return
		}
		log.Println("call error ", panicErr)
		err = errors.New("invoke failed!")
	}()
	inParam, err := this.decodeInParam(inBytes)
	if err != nil {
		return nil, err
	}

	paramV := []reflect.Value{reflect.ValueOf(inParam)}
	retV := this.fun.Call(paramV)
	if !retV[1].IsNil() {
		err = retV[1].Interface().(error)
	}
	outParam := retV[0].Interface()
	outBytes, err = std.MsgpackMarshal(outParam)
	if err != nil {
		return nil, err
	}
	return outBytes, nil
}

type RPC struct {
	ioLoop       *liblpc.IOEvtLoop
	rcpFuncMap   map[string]*rpcFunc
	promiseGroup *std.PromiseGroup
	lock         *sync.RWMutex
	startFlag    int32
}

func New() (*RPC, error) {
	loop, err := liblpc.NewIOEvtLoop(1024 * 1024 * 4)
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

func checkInParam(t reflect.Type) {
	inNum := t.NumIn()
	std.Assert(inNum == 1, "func in param len != 1")
	in := t.In(0)
	inKind := in.Kind()
	std.Assert(inKind == reflect.Ptr || inKind == reflect.Struct, "param must be prt of struct")
}

func checkOutParam(t reflect.Type) {
	outNum := t.NumOut()
	std.Assert(outNum == 2, "func out param len != 2 ")
	out1 := t.Out(1)
	std.Assert(out1 == typeOfError, "out1 must be error type")
}

func getFuncName(fv reflect.Value) string {
	fname := runtime.FuncForPC(reflect.Indirect(fv).Pointer()).Name()
	idx := strings.LastIndex(fname, ".")
	return fname[idx+1:]
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

func getValueElement(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

func (this *RPC) Call(sw liblpc.StreamWriter, timeout time.Duration, name string, param interface{}, out interface{}) error {
	outMsg := &rpcRawMsg{
		Id:         std.GenRandomUUID(),
		MethodName: name,
		Type:       rpcReqMsg,
	}
	err := outMsg.SetData(param)
	if err != nil {
		return err
	}
	//
	outBytes, err := encodeRpcMsg(outMsg)
	sw.Write(outBytes, false)
	//
	promise := std.NewPromise()
	promiseId := std.PromiseId(outMsg.Id)
	this.promiseGroup.AddPromise(promiseId, promise)
	defer this.promiseGroup.RemovePromise(promiseId)
	future := promise.GetFuture()
	data, err := future.WaitData(timeout)
	if err != nil {
		return err
	}
	dataBytes, ok := data.([]byte)
	std.Assert(ok, "data not bytes!")
	err = std.MsgpackUnmarshal(dataBytes, out)
	if err != nil {
		return err
	}
	return nil
}

func (this *RPC) Start() {
	if atomic.CompareAndSwapInt32(&this.startFlag, 0, 1) {
		go this.ioLoop.Run()
	}
}

func (this *RPC) Close() error {
	return this.ioLoop.Close()
}

func (this *RPC) AddApiStream(fd int, userData interface{}) liblpc.StreamWriter {
	s := liblpc.NewFdBufferedStream(this.ioLoop, fd, this.genericRead)
	s.SetUserData(userData)
	s.Start()
	return s
}

const kMaxRpcMsgBodyLen = 1024 * 1024 * 32

func (this *RPC) genericRead(sw liblpc.StreamWriter, buf std.ReadableBuffer, err error) {
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
