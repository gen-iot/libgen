package rpc

import (
	"fmt"
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/liblpc"
	"testing"
	"time"
)

type Req struct {
	A  int
	B  int
	Tm time.Time
}

type Rsp struct {
	Sum int
}

func sum(req Req) (Rsp, error) {
	fmt.Println("req delta time -> ", time.Now().Sub(req.Tm))
	return Rsp{
		Sum: req.A + req.B,
	}, nil
}

func startLocalRpc(fd int) {
	rpc, err := New()
	std.AssertError(err, "new rpc")
	rpc.AddApiStream(fd, nil)
	rpc.RegFun(sum)
	rpc.Start()
}

func startMockRemoteRpc(fd int) {
	rpc, err := New()
	std.AssertError(err, "new rpc")
	sw := rpc.AddApiStream(fd, nil)
	rpc.Start()
	out := new(Rsp)
	err = rpc.Call(sw, time.Second*60, "sum", &Req{
		A:  10,
		B:  100,
		Tm: time.Now(),
	}, out)
	std.AssertError(err, "call error")
	std.Assert(out.Sum == 10+100, "result error")
}

func TestRpc(t *testing.T) {
	fds, err := liblpc.MakeIpcSockpair(true)
	std.AssertError(err, "socketPair error")
	startLocalRpc(fds[0])
	startMockRemoteRpc(fds[1])
}
