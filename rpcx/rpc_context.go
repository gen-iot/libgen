package rpcx

type Context interface {
	Callable() Callable

	SetMethod(string)
	Method() string

	SetHeaders(h map[string]string)
	Headers() map[string]string

	SetRequest(in interface{})
	Request() interface{}

	SetResponse(out interface{})
	Response() interface{}

	SetError(err error)
	Error() error

	Direction() Direction

	setDirection(d Direction)
}

type contextImpl struct {
	call      Callable
	method    string
	headers   map[string]string
	in        interface{}
	out       interface{}
	err       error
	direction Direction
	rawInMsg  *rpcRawMsg
	rawOutMsg *rpcRawMsg
}

func (this *contextImpl) Method() string {
	return this.method
}

func (this *contextImpl) SetMethod(method string) {
	this.method = method
}

func (this *contextImpl) Headers() map[string]string {
	return this.headers
}

func (this *contextImpl) SetHeaders(h map[string]string) {
	this.headers = h
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

func newContext(call Callable, inMsg *rpcRawMsg) *contextImpl {
	return &contextImpl{
		call:      call,
		method:    inMsg.MethodName,
		headers:   inMsg.Headers,
		direction: In,
		rawInMsg:  inMsg,
	}
}
