package rpcx

import (
	"fmt"
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/libgen"
	"gitee.com/SuzhenProjects/liblpc"
	"net"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"syscall"
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

func sum(call Callable, req *Req) (*Rsp, error) {
	fmt.Println("req delta time -> ", time.Now().Sub(req.Tm))
	return &Rsp{
		Sum: req.A + req.B,
	}, nil
}

func startLocalRpc(fd int, wg *sync.WaitGroup) {
	rpc, err := New()
	std.AssertError(err, "new rpcx")
	defer std.CloseIgnoreErr(rpc)
	rpc.Start()
	rpc.RegFun(sum)
	rpc.NewCallable(fd, nil)
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
		err = callable.Call(time.Second*5, "sum", &Req{
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
	//
	fds, err := liblpc.MakeIpcSockpair(true)
	std.AssertError(err, "socketPair error")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startMockRemoteRpc(fds[1], wg)
	startLocalRpc(fds[0], wg)
}

func ____Ping(callable Callable, req *libgen.Ping) (*libgen.Pong, error) {
	fmt.Println("recv ping ->", req, ", delta is = ", time.Now().Sub(req.Time))
	return &libgen.Pong{
		Time: time.Now(), Msg: "KERNEL PONG",
	}, nil
}

func TestRemoteTcpRpc(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8000")
	std.AssertError(err, "ResolveTCPAddr")
	listener, err := net.ListenTCP("tcp", addr)
	std.AssertError(err, "ListenTCP")
	conn, err := listener.AcceptTCP()
	std.AssertError(err, "AcceptTCP")
	fmt.Println("NewConnection")
	_ = conn.SetNoDelay(true)
	f, err := conn.File()
	std.AssertError(err, "getFd")
	fd := int(f.Fd())
	rpc, err := New()
	std.AssertError(err, "new rpc failed")
	rpc.Start()
	rpc.RegFuncWithName("Ping", ____Ping)
	callable := rpc.NewCallable(fd, nil)
	fmt.Println("NewCallable")
	clientRsp := new(libgen.Pong)
	for {
		runtime.KeepAlive(conn)
		time.Sleep(500 * time.Millisecond)
		err = callable.Call(time.Second, "Ping", &libgen.Ping{
			Time: time.Now(),
			Msg:  "ping from server",
		}, clientRsp)
		std.AssertError(err, "Call Ping")
		fmt.Println("client resp -> ", clientRsp)
	}
}

func TestRemoteTcpRpcV2(t *testing.T) {
	listenFd, err := syscall.Socket(syscall.AF_INET, syscall.SOL_SOCKET|syscall.SOCK_CLOEXEC, syscall.IPPROTO_TCP)
	std.AssertError(err, "create listen socket failed")
	defer func() {
		_ = syscall.Close(listenFd)
	}()
	err = syscall.Bind(listenFd, &syscall.SockaddrInet4{
		Port: 8000,
		Addr: [4]byte{0, 0, 0, 0},
	})
	std.AssertError(err, "bind err")
	err = syscall.Listen(listenFd, 128)
	std.AssertError(err, "listen err")
	nfd, _, err := syscall.Accept4(listenFd, syscall.O_NONBLOCK|syscall.O_CLOEXEC)
	std.AssertError(err, "accept err")
	rpc, err := New()
	std.AssertError(err, "new rpc failed")
	rpc.Start()
	rpc.RegFuncWithName("Ping", ____Ping)
	callable := rpc.NewCallable(nfd, nil)
	fmt.Println("NewCallable")
	clientRsp := new(libgen.Pong)
	for {
		time.Sleep(5 * time.Millisecond)
		err = callable.Call(time.Second, "Ping", &libgen.Ping{
			Time: time.Now(),
			Msg:  "ping from server",
		}, clientRsp)
		std.AssertError(err, "Call Ping")
		fmt.Println("client resp -> ", clientRsp)
	}
}
