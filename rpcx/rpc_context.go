package rpcx

type Context interface {
	Callable() Callable

	Id() string

	SetMethod(string)
	Method() string

	SetHeaders(d Direction, h map[string]string)
	Headers(d Direction) map[string]string

	SetRequest(in interface{})
	Request() interface{}

	SetResponse(out interface{})
	Response() interface{}

	SetError(err error)
	Error() error

	Direction() Direction
}

type contextImpl struct {
	call      Callable
	in        interface{}
	out       interface{}
	err       error
	direction Direction
	inMsg     *rpcRawMsg
	outHeader map[string]string
}

func (this *contextImpl) Method() string {
	return this.inMsg.MethodName
}

func (this *contextImpl) SetMethod(method string) {
	this.inMsg.MethodName = method
}

func (this *contextImpl) Headers(d Direction) map[string]string {
	if d == In {
		return this.inMsg.Headers
	} else if d == Out {
		return this.outHeader
	}
	return nil
}

func (this *contextImpl) SetHeaders(d Direction, h map[string]string) {
	if d == In {
		this.inMsg.Headers = h
	} else if d == Out {
		this.outHeader = h
	}
}

func (this *contextImpl) Id() string {
	return this.inMsg.Id
}

func (this *contextImpl) SetRequest(in interface{}) {
	this.in = in
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

func (this *contextImpl) Direction() Direction {
	return this.direction
}

func (this *contextImpl) setDirection(d Direction) {
	this.direction = d
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
		call:      call,
		direction: In,
		inMsg:     inMsg,
	}
}
