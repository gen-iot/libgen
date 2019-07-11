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
	mid []MiddlewareFunc
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
	outMsg := &rpcRawMsg{
		Id:         std.GenRandomUUID(),
		MethodName: name,
		Headers:    headers,
		Type:       rpcReqMsg,
	}
	err := outMsg.SetData(param)
	if err != nil {
		return err
	}
	// log.Println("SEND REQ id -> ", outMsg.Id)
	//add promise
	promise := std.NewPromise()
	promiseId := std.PromiseId(outMsg.Id)
	this.ctx.promiseGroup.AddPromise(promiseId, promise)
	defer this.ctx.promiseGroup.RemovePromise(promiseId)
	//write out
	outBytes, err := encodeRpcMsg(outMsg)
	if err != nil {
		return err
	}
	this.stream.Write(outBytes, false)
	//wait for data
	future := promise.GetFuture()
	data, err := future.WaitData(timeout)
	if err != nil {
		log.Println("Call :future wait got err ->", err)
		return err
	}
	dataBytes, ok := data.([]byte)
	std.Assert(ok, "Call :data not bytes!")
	err = std.MsgpackUnmarshal(dataBytes, out)
	if err != nil {
		log.Println("Call :MsgpackUnmarshal got err ->", err)
		return err
	}
	return nil
}
