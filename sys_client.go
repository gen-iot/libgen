//+build client

package libgen

import (
	"fmt"
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/libgen/rpcx"
	"log"
	"net"
	"runtime"
	"sync"
	"time"
)

var initOnce = sync.Once{}
var gCallable rpcx.Callable
var gRpc *rpcx.RPC
var gApiClient *ApiClientImpl

var ApiCallTimeout = time.Second * 30

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
	addr, err := net.ResolveTCPAddr("tcp", "192.168.50.232:8000")
	std.AssertError(err, "ResolveTCPAddr error")
	conn, err := net.DialTCP("tcp", nil, addr)
	runtime.SetFinalizer(conn, func(conn *net.TCPConn) {
		fmt.Println("close conn")
	})
	runtime.KeepAlive(conn)
	std.AssertError(err, "net dial error")
	err = conn.SetNoDelay(true)
	std.AssertError(err, "net SetNoDelay error")
	f, err := conn.File()
	std.AssertError(err, "get conn file error")
	fd := int(f.Fd())

	gCallable = gRpc.NewCallable(fd, nil)

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
