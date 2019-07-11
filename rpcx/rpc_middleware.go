package rpcx

type Direction = int

const (
	In Direction = iota
	Out
)

type HandleFunc = func(ctx Context) error

type MiddlewareFunc = func(next HandleFunc) HandleFunc
