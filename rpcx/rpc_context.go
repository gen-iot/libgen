package rpcx

import "gitee.com/SuzhenProjects/liblpc"

type Context interface {
	Callable() Callable

	Id() string

	SetMethod(string)
	Method() string

	SetRequest(in interface{})
	Request() interface{}

	SetResponse(out interface{})
	Response() interface{}

	SetError(err error)
	Error() error

	liblpc.UserDataStorage
}

type contextImpl struct {
	call   Callable
	in     interface{}
	out    interface{}
	err    error
	reqMsg *rpcRawMsg
	ackMsg *rpcRawMsg
	liblpc.BaseUserData
}

func (this *contextImpl) Method() string {
	return this.reqMsg.MethodName
}

func (this *contextImpl) SetMethod(method string) {
	this.reqMsg.MethodName = method
}

func (this *contextImpl) Id() string {
	return this.reqMsg.Id
}

func (this *contextImpl) SetRequest(in interface{}) {
	this.in = in
	_ = this.reqMsg.SetData(in)
}

func (this *contextImpl) Request() interface{} {
	return this.in
}

func (this *contextImpl) SetResponse(out interface{}) {
	this.out = out
}

func (this *contextImpl) Response() interface{} {
	return this.out
}

func (this *contextImpl) SetError(err error) {
	this.err = err
}

func (this *contextImpl) Error() error {
	return this.err
}

func (this *contextImpl) Callable() Callable {
	return this.call
}

func (this *contextImpl) buildOutMsg() *rpcRawMsg {
	out := &rpcRawMsg{
		Id:         this.Id(),
		MethodName: this.Method(),
		Type:       rpcAckMsg,
	}
	if this.err != nil {
		out.SetError(this.err)
	} else {
		err := out.SetData(this.out)
		if err != nil {
			out.SetError(err)
		}
	}
	return out
}

func newContext(call Callable, inMsg *rpcRawMsg) *contextImpl {
	return &contextImpl{
		call:   call,
		reqMsg: inMsg,
	}
}
