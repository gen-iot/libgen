//+build !server

package libgen

import (
	"fmt"
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/libgen/rpcx"
	"gitee.com/SuzhenProjects/liblpc"
	"sync"
	"time"
)

var initOnce = sync.Once{}
var gCallable rpcx.Callable
var gRpc *rpcx.RPC
var gApiClient *ApiClientImpl

var ApiCallTimeout = time.Second * 1

const clientFd = 3

type Config struct {
	//1 local ;2 remote
	Mode          int
	RemoteAddress string
}

var DefaultConfig = Config{

}

func Init() {
	InitWithConfig(DefaultConfig)
}

func InitWithConfig(config Config) {
	initOnce.Do(func() {
		doInit(config)
	})
}

func doInit(config Config) {
	fmt.Println("LIBGEN CLIENT INIT")
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpc failed")
	gRpc = rpc
	//todo 根据Manifest决定是否注册
	gRpc.RegFuncWithName("", onDeviceControl)
	gRpc.RegFuncWithName("", onDeviceStatus)
	gRpc.RegFuncWithName("Ping", onPing)
	gRpc.Start()
	if config.Mode == 2 {
		sockFd, err := liblpc.NewConnFdSimple(config.RemoteAddress)
		std.AssertError(err, "connect err")
		gCallable = gRpc.NewCallable(int(sockFd), nil)
	} else {
		gCallable = gRpc.NewCallable(clientFd, nil)
	}

	gApiClient = new(ApiClientImpl)
	fmt.Println("LIBGEN CLIENT INIT SUCCESS")
}

func getCallable() rpcx.Callable {
	return gCallable
}

func GetApiClient() *ApiClientImpl {
	return gApiClient
}