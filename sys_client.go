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
	Mode          int     `json:"mode" validate:"required,oneof=1 2"`
	RemoteAddress string  `json:"remoteAddress"`
	PkgInfo       PkgInfo `json:"pkgInfo" validate:"required"`
	AccessToken   string  `json:"accessToken" validate:"required"`
}

var DefaultConfig = Config{
	Mode:          1,
	RemoteAddress: "",
	AccessToken:   "",
	PkgInfo: PkgInfo{
		Package: "",
		Name:    "",
	},
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
	err := std.ValidateStruct(config)
	std.AssertError(err, "config invalid")
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpc failed")
	gRpc = rpc
	//todo 根据Manifest决定是否注册
	//gRpc.RegFuncWithName("", onDeviceControl)
	//gRpc.RegFuncWithName("", onDeviceStatus)
	gRpc.RegFuncWithName("Ping", onPing)
	gRpc.Start()
	if config.Mode == 2 {
		sockFd, err := liblpc.NewConnFd(config.RemoteAddress)
		std.AssertError(err, "connect err")
		gCallable = gRpc.NewCallable(int(sockFd), nil)
		//handshake
		out := new(BaseResponse)
		err = gCallable.Call(ApiCallTimeout, "Handshake", &HandshakeRequest{
			PkgInfo:     config.PkgInfo,
			AccessToken: config.AccessToken,
		}, out)
		std.AssertError(err, "connect handshake err")
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
