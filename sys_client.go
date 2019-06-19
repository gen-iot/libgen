//+build client

package libgen

import (
	"fmt"
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/libgen/rpcx"
	"log"
	"sync"
	"syscall"
	"time"
)

var initOnce = sync.Once{}
var gCallable rpcx.Callable
var gRpc *rpcx.RPC
var gApiClient *ApiClientImpl

var ApiCallTimeout = time.Second * 1

const clientFd = 3

func Init() {
	initOnce.Do(doInit)
}

func doInit() {
	fmt.Println("LIBGEN CLIENT INIT")
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpc failed")
	gRpc = rpc
	//gRpc.RegFun(deviceControl)
	gRpc.RegFuncWithName("Ping", onPing)
	gRpc.Start()
	gApiClient = new(ApiClientImpl)
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOL_SOCKET, syscall.IPPROTO_TCP)
	std.AssertError(err, "new sock err")
	err = syscall.Connect(sock, &syscall.SockaddrInet4{
		Port: 8000,
		Addr: [4]byte{192, 168, 50, 48},
	})
	std.AssertError(err, "connect err")
	gCallable = gRpc.NewCallable(sock, nil)
	fmt.Println("LIBGEN CLIENT INIT SUCCESS")
}

func getCallable() rpcx.Callable {
	return gCallable
}

func GetApiClient() *ApiClientImpl {
	return gApiClient
}

var gOnDeviceControl OnDeviceControlFun

type OnDeviceControlFun func(req *ControlDeviceRequest) (*ControlDeviceResponse, error)
type OnDeviceStatusFun func(req *ControlDeviceRequest) (*ControlDeviceResponse, error)

func RegOnDeviceControl(fn OnDeviceControlFun) {
	gOnDeviceControl = fn
}

func RegOnDeviceStatus(fn OnDeviceControlFun) {
	gOnDeviceControl = fn
}

func onPing(callable rpcx.Callable, req *Ping) (*Pong, error) {
	log.Println("receive ping req.time =", req.Time, " delta is ", time.Now().Sub(req.Time))
	return &Pong{Time: time.Now(), Msg: "client pong"}, nil
}

func onDeviceControl(req *ControlDeviceRequest) (*ControlDeviceResponse, error) {
	if gOnDeviceControl != nil {
		return gOnDeviceControl(req)
	}
	return new(ControlDeviceResponse), nil
}

func onDeviceStatus(req *ControlDeviceRequest) (*ControlDeviceResponse, error) {
	if gOnDeviceControl != nil {
		return gOnDeviceControl(req)
	}
	return new(ControlDeviceResponse), nil
}
