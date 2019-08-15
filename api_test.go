package libgen

import (
	"fmt"
	"github.com/gen-iot/liblpc"
	"github.com/gen-iot/rpcx"
	"github.com/gen-iot/std"
	"log"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestApiClientImpl_Ping(t *testing.T) {
	wg := sync.WaitGroup{}
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpc failed")
	//gRpc.RegFunc(deviceControl)
	rpc.RegFuncWithName("Ping", pong)
	rpc.Start()
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOL_SOCKET, syscall.IPPROTO_TCP)
	std.AssertError(err, "new sock err")
	err = syscall.Connect(sock, &syscall.SockaddrInet4{
		Port: 8000,
		Addr: [4]byte{192, 168, 50, 48},
	})
	std.AssertError(err, "connect err")
	callable := rpc.NewConnCallable(sock, nil)
	callable.Start()
	wg.Add(1)
	go func() {
		defer wg.Done()
		count := 0
		for {
			log.Println("ping test count :", count)
			res := new(Pong)
			err = callable.Call5(ApiCallTimeout, "Ping", &Ping{Time: time.Now(), Msg: fmt.Sprintf("client ping %d", count)}, res)
			std.AssertError(err, "ping error")
			log.Println("ping res msg >> ", res.Msg)
			time.Sleep(time.Millisecond * 1)
			count++
			//if count >= 10000 {
			//	break
			//}
		}
		//time.Sleep(time.Second * 60 * 60)
	}()
	wg.Wait()
	log.Println("ping test over")
}

const remoteAddr = "192.168.50.232:54321"

func TestApiClientImpl_ControlDevice(t *testing.T) {
	pkg := PkgInfo{
		Package: "com.gen.kernel",
		Name:    "Manage",
	}
	sockFd, err := liblpc.NewConnFd(remoteAddr)
	std.AssertError(err, "connect err")
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpc")
	defer std.CloseIgnoreErr(rpc)
	rpc.Start()
	callable := rpc.NewConnCallable(int(sockFd), nil)
	//handshake
	err = callable.Call1(time.Second*10, "Handshake", &HandshakeRequest{
		PkgInfo:     pkg,
		AccessToken: "pujie123",
	})
	std.AssertError(err, "Handshake failed")
	err = callable.Call1(time.Second*10, "ControlDevice", &ControlDeviceRequest{
		PkgInfo: PkgInfo{
			Package: "com.pujie88.iot",
			Name:    "HotelRemote",
		},
		Id: "014100000000936A_0_0_67",
		CtrlParams: map[string]interface{}{
			"power": 1,
		},
	})
	std.AssertError(err, "control failed")
}

func TestDupSocket(t *testing.T) {
	fd1, err := liblpc.NewTcpSocketFd(4, false, true)
	std.AssertError(err, "new sock")
	baiduAddr, err := liblpc.ResolveTcpAddr("www.baidu.com:80")
	std.AssertError(err, "resolve baidu addr err")
	err = syscall.Connect(int(fd1), baiduAddr)
	std.AssertError(err, "connect err")
	fd2, err := liblpc.NewTcpSocketFd(4, true, true)
	std.AssertError(err, "new sock2")
	err = syscall.Dup2(int(fd1), int(fd2))
	std.AssertError(err, "dup sock err")
}
