package libgen

import (
	"fmt"
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/libgen/rpcx"
	"log"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestApiClientImpl_Ping(t *testing.T) {
	Init()
	wg := sync.WaitGroup{}
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpc failed")
	//gRpc.RegFun(deviceControl)
	rpc.RegFuncWithName("Ping", onPing)
	rpc.Start()
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOL_SOCKET, syscall.IPPROTO_TCP)
	std.AssertError(err, "new sock err")
	err = syscall.Connect(sock, &syscall.SockaddrInet4{
		Port: 8000,
		Addr: [4]byte{192, 168, 50, 48},
	})
	std.AssertError(err, "connect err")
	callable := rpc.NewCallable(sock, nil)
	wg.Add(1)
	go func() {
		count := 0
		for {
			log.Println("ping test count :", count)
			res := new(Pong)
			err=callable.Call(ApiCallTimeout, "Ping", &Ping{Time: time.Now(), Msg: fmt.Sprintf("client ping %d", count)}, res)
			std.AssertError(err, "ping error")
			log.Println("ping res msg >> ", res.Msg)
			time.Sleep(time.Millisecond * 1)
			count++
			//if count >= 10000 {
			//	break
			//}
		}
		//time.Sleep(time.Second * 60 * 60)
		wg.Done()
	}()
	wg.Wait()
	log.Println("ping test over")
}
