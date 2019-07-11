package rpcx

import "gitee.com/Puietel/std"

type Direction = int

const (
	In Direction = iota
	Out
)

type HandleFunc = func(ctx Context) error

var emptyHandlerFunc HandleFunc = func(ctx Context) error {
	return nil
}

type MiddlewareFunc = func(next HandleFunc) HandleFunc

type middleware struct {
	midwares []MiddlewareFunc
}

func (this *middleware) Use(m ...MiddlewareFunc) {
	if this.midwares == nil {
		this.midwares = make([]MiddlewareFunc, 0, 4)
	}
	this.midwares = append(this.midwares, m...)
}

var emptyHandleFunc = func(ctx Context) error {
	return nil
}

func (this *middleware) buildChain(direct Direction, h HandleFunc) HandleFunc {
	std.Assert(h != nil, "buildMiddleware, h == nil")
	switch direct {
	case In:
		{
			for i := len(this.midwares) - 1; i >= 0; i-- {
				h = this.midwares[i](h)
			}
		}
	case Out:
		{
			xh := emptyHandleFunc
			for i := 0; i < len(this.midwares); i++ {
				xh = this.midwares[i](xh)
			}
			return func(ctx Context) error {
				_ = h(ctx)
				return xh(ctx)
			}

		}
	default:
		std.Assert(false, "unknown direction")
	}
	return h
}
