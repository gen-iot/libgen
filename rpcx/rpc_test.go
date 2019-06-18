package rpcx

import (
	"fmt"
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/liblpc"
	"os"
	"runtime/pprof"
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
	std.AssertError(err, "new rpcx")
	defer std.CloseIgnoreErr(rpc)
	rpc.NewCallable(fd, nil)
	rpc.RegFun(sum)
	rpc.Start()
	wg.Wait()
}

func startMockRemoteRpc(fd int, wg *sync.WaitGroup) {
	defer wg.Done()
	rpc, err := New()
	std.AssertError(err, "new rpcx")
	defer std.CloseIgnoreErr(rpc)
	rpc.Start()
	callable := rpc.NewCallable(fd, nil)
	after := time.After(time.Second * 5)
	for {
		select {
		case <-after:
			return
		default:
		}
		out := new(Rsp)
		err = callable.Call(time.Second, "sum", &Req{
			A:  10,
			B:  100,
			Tm: time.Now(),
		}, out)
		std.AssertError(err, "call error")
		std.Assert(out.Sum == 10+100, "result error")
	}
}

func TestRpc(t *testing.T) {
	file, err := os.OpenFile("cpu_prof.prof", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	std.AssertError(err, "create prof failed")
	defer std.CloseIgnoreErr(file)
	err = pprof.StartCPUProfile(file)
	std.AssertError(err, "start profile failed")
	defer pprof.StopCPUProfile()
	fds, err := liblpc.MakeIpcSockpair(true)
	std.AssertError(err, "socketPair error")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startMockRemoteRpc(fds[1], wg)
	startLocalRpc(fds[0], wg)
}