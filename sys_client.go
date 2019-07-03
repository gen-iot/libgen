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

var ApiCallTimeout = time.Second * 5

const clientFd = 3

type config struct {
	Type        AppType `json:"type" validate:"required,oneof=900 901"`
	Endpoint    string  `json:"endpoint"`
	PkgInfo     PkgInfo `json:"pkgInfo" validate:"required"`
	AccessToken string  `json:"accessToken" validate:"required"`
}

var defaultConfig = config{
	Type:        LocalApp,
	Endpoint:    "",
	AccessToken: "",
	PkgInfo: PkgInfo{
		Package: "",
		Name:    "",
	},
}

func InitLocal() {
	initWithConfig(defaultConfig)
}

func InitRemote(endPoint string, pkgInfo PkgInfo, accessToken string) {
	initWithConfig(config{
		Type:        RemoteApp,
		Endpoint:    endPoint,
		PkgInfo:     pkgInfo,
		AccessToken: accessToken,
	})
}

func initWithConfig(config config) {
	initOnce.Do(func() {
		doInit(config)
	})
}

func Cleanup() {
	std.CloseIgnoreErr(gCallable)
	std.CloseIgnoreErr(gRpc)
	gApiClient = nil
}

func doInit(config config) {
	fmt.Println("LIBGEN CLIENT INIT")
	err := std.ValidateStruct(config)
	std.AssertError(err, "config invalid")
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpc failed")
	gRpc = rpc
	gRpc.RegFuncWithName("ControlDevice", onDeviceControl)
	gRpc.RegFuncWithName("DeliveryDeviceStatus", onDeviceStatusDelivery)
	gRpc.RegFuncWithName("Ping", pong)
	gRpc.OnCallableClosed(onCallableClose)
	gRpc.Start()
	if config.Type == RemoteApp {
		sockFd, err := liblpc.NewConnFd(config.Endpoint)
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

func onCallableClose(callable rpcx.Callable) {
	fmt.Println("LIBGEN RPC DISCONNECTED ")
	//todo 重连

}

func getCallable() rpcx.Callable {
	return gCallable
}

func GetApiClient() RpcApiClient {
	std.Assert(gApiClient != nil, "please init first")
	return gApiClient
}
