package rpcx

type Context interface {
	Callable() Callable

	Id() string

	SetMethod(string)
	Method() string

	RequestBytes() []byte

	SetRequest(in interface{})
	Request() interface{}

	SetResponse(out interface{})
	Response() interface{}

	SetError(err error)
	Error() error
}

type contextImpl struct {
	call      Callable
	in        interface{}
	out       interface{}
	err       error
	inMsg     *rpcRawMsg
	outHeader map[string]string
}

func (this *contextImpl) Method() string {
	return this.inMsg.MethodName
}

func (this *contextImpl) SetMethod(method string) {
	this.inMsg.MethodName = method
}

func (this *contextImpl) Id() string {
	return this.inMsg.Id
}

func (this *contextImpl) SetRequest(in interface{}) {
	this.in = in
}

func (this *contextImpl) RequestBytes() []byte {
	return this.inMsg.Data
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
		call:  call,
		inMsg: inMsg,
	}
}
