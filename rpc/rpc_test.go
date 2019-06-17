package rpc

import (
	"fmt"
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/liblpc"
	"sync"
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

func sum(req *Req) (*Rsp, error) {
	fmt.Println("req delta time -> ", time.Now().Sub(req.Tm))
	return &Rsp{
		Sum: req.A + req.B,
	}, nil
}

func startLocalRpc(fd int, wg *sync.WaitGroup) {
	rpc, err := New()
	std.AssertError(err, "new rpc")
	defer std.CloseIgnoreErr(rpc)
	sw := rpc.NewCallable(fd, nil)
	defer std.CloseIgnoreErr(sw)
	rpc.RegFun(sum)
	rpc.Start()
	wg.Wait()
}

func startMockRemoteRpc(fd int, wg *sync.WaitGroup) {
	defer wg.Done()
	rpc, err := New()
	std.AssertError(err, "new rpc")
	defer std.CloseIgnoreErr(rpc)
	callable := rpc.NewCallable(fd, nil)
	rpc.Start()
	after := time.After(time.Second * 10)
	for {
		select {
		case <-after:
			return
		default:
		}
		out := new(Rsp)
		err = callable.Call(time.Second*1, "sum", &Req{
			A:  10,
			B:  100,
			Tm: time.Now(),
		}, out)
		std.AssertError(err, "call error")
		std.Assert(out.Sum == 10+100, "result error")
		time.Sleep(time.Millisecond * 500)
	}
}

func TestRpc(t *testing.T) {
	fds, err := liblpc.MakeIpcSockpair(true)
	std.AssertError(err, "socketPair error")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startMockRemoteRpc(fds[1], wg)
	startLocalRpc(fds[0], wg)
}
