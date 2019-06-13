package libgen

import (
	"fmt"
	"gitee.com/Puietel/std"
	"testing"
)

type Req struct {
	A int
	B int
}

type Rsp struct {
	Sum int
}

func sum(req Req) (Rsp, error) {
	return Rsp{
		Sum: req.A + req.B,
	}, nil
}

func TestRegFn(t *testing.T) {
	rpc := NewRpc()


	rpc.RegFun(sum)


	reply := new(Rsp)

	err := rpc.MockCall("sum", Req{
		A: 10,
		B: 20,
	}, reply)


	std.Assert(err == nil, "sum")
	fmt.Println("reply -> ", *reply)
}
