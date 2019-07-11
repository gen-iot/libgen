package libgen

import (
	"fmt"
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/libgen/rpcx"
	"syscall"
	"testing"
	"time"
)

func ____Ping(ctx rpcx.Context, req *Ping) (*Pong, error) {
	fmt.Println("recv ping ->", req, ", delta is = ", time.Now().Sub(req.Time))
	return &Pong{
		Time: time.Now(), Msg: "KERNEL PONG",
	}, nil
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
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpc failed")
	rpc.Start()
	rpc.RegFuncWithName("Ping", ____Ping)
	callable := rpc.NewConnCallable(nfd, nil)
	fmt.Println("NewConnCallable")
	clientRsp := new(Pong)
	for {
		time.Sleep(5 * time.Millisecond)
		err = callable.Call(time.Second, "Ping", &Ping{
			Time: time.Now(),
			Msg:  "ping from server",
		}, clientRsp)
		std.AssertError(err, "call Ping")
		fmt.Println("client resp -> ", clientRsp)
	}
}
