//+build client

package libgen

import (
	"fmt"
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/libgen/rpcx"
	"sync"
	"time"
)

var initOnce = sync.Once{}
var gCallable rpcx.Callable
var gRpc *rpcx.RPC
var gApiClient *ApiClientImpl

var ApiCallTimeout = time.Second * 5

const clientFd = 3

func Init() {
	initOnce.Do(doInit)
}

func doInit() {
	fmt.Println("LIBGEN CLIENT INIT")
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpc failed")
	gRpc = rpc
	gRpc.RegFun(deviceControl)
	gRpc.RegFun(deviceStatus)
	gRpc.Start()
	gApiClient = new(ApiClientImpl)
	gCallable = gRpc.NewCallable(clientFd, nil)
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

func deviceControl(req *ControlDeviceRequest) (*ControlDeviceResponse, error) {
	if gOnDeviceControl != nil {
		return gOnDeviceControl(req)
	}
	return new(ControlDeviceResponse), nil
}
func deviceStatus(req *ControlDeviceRequest) (*ControlDeviceResponse, error) {
	if gOnDeviceControl != nil {
		return gOnDeviceControl(req)
	}
	return new(ControlDeviceResponse), nil
}
