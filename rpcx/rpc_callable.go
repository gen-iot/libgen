package rpcx

import (
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/liblpc"
	"io"
	"log"
	"time"
)

type Callable interface {
	io.Closer
	liblpc.UserDataStorage
	Call(timeout time.Duration, name string, param interface{}, out interface{}) error
	CallWithHeader(timeout time.Duration, name string, headers map[string]string, param interface{}, out interface{}) error
}

type rpcCli struct {
	stream *liblpc.BufferedStream
	ctx    *RPC
	middleware
	liblpc.BaseUserData
}

func (this *rpcCli) start() {
	this.stream.Start()
}

func (this *rpcCli) Close() error {
	return this.stream.Close()
}

func (this *rpcCli) Call(timeout time.Duration, name string, param interface{}, out interface{}) error {
	return this.CallWithHeader(timeout, name, nil, param, out)
}

func (this *rpcCli) CallWithHeader(timeout time.Duration, name string, headers map[string]string, param interface{}, out interface{}) error {
	std.Assert(this.stream != nil, "stream is nil!")
	msgId := std.GenRandomUUID()
	msg := &rpcRawMsg{
		Id:         msgId,
		MethodName: name,
		Headers:    headers,
		Type:       rpcReqMsg,
	}
	//add promise
	ctx := newContext(this, msg)
	ctx.SetRequest(param)
	f := this.buildInvoke(timeout, ctx, out)
	h := this.buildChain(f)
	h(ctx)
	return ctx.Error()
}

func (this *rpcCli) buildInvoke(timeout time.Duration, ctx *contextImpl, out interface{}) HandleFunc {
	return func(Context) {
		this.invoke(timeout, out, ctx)
	}
}

func (this *rpcCli) invoke(timeout time.Duration, out interface{}, ctx *contextImpl) {
	err := ctx.inMsg.SetData(ctx.in)
	if err != nil {
		ctx.SetError(err)
		return
	}
	promise := std.NewPromise()
	promiseId := std.PromiseId(ctx.Id())
	//write out
	outBytes, err := encodeRpcMsg(ctx.inMsg)
	if err != nil {
		ctx.SetError(err)
		return
	}
	this.ctx.promiseGroup.AddPromise(promiseId, promise)
	defer this.ctx.promiseGroup.RemovePromise(promiseId)
	//
	this.stream.Write(outBytes, false)
	//wait for data
	future := promise.GetFuture()
	data, err := future.WaitData(timeout)
	if err != nil {
		log.Println("call :future wait got err ->", err)
		ctx.SetError(err)
		return
	}
	dataBytes, ok := data.([]byte)
	std.Assert(ok, "call :data not bytes!")
	err = std.MsgpackUnmarshal(dataBytes, out)
	if err != nil {
		log.Println("call :MsgpackUnmarshal got err ->", err)
		ctx.SetError(err)
	}
	ctx.SetResponse(out)
}
