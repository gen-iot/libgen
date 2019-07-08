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
var gRpc *rpcx.RPC
var gApiClient *ApiClientImpl
var gConfig = defaultConfig
var gOnConnected func() = nil

var ApiCallTimeout = time.Second * 5

const clientFd = 3

type config struct {
	Type        AppType `json:"type" validate:"required,oneof=900 901"`
	Endpoint    string  `json:"endpoint"`
	PkgInfo     PkgInfo `json:"pkgInfo" validate:"-"`
	AccessToken string  `json:"accessToken"`
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

func InitLocal(onConnected func()) {
	initWithConfig(defaultConfig)
	gOnConnected = onConnected
	connect()
}

func InitRemote(endPoint string, pkgInfo PkgInfo, accessToken string, onConnected func()) {
	initWithConfig(config{
		Type:        RemoteApp,
		Endpoint:    endPoint,
		PkgInfo:     pkgInfo,
		AccessToken: accessToken,
	})
	gOnConnected = onConnected
	connect()
}

func newRemoteCallable(endpoint string, timeout time.Duration) (rpcx.Callable, error) {
	sockFd, err := liblpc.NewConnFd(endpoint)
	if err != nil {
		return nil, err
	}
	cancelFn, future := gRpc.NewClientCallable(int(sockFd), nil)
	// sync wait
	data, err := future.WaitData(timeout)
	if err != nil {
		cancelFn()
		return nil, err
	}
	return data.(rpcx.Callable), nil
}

func connect() {
	std.Assert(gConfig.Type == LocalApp || gConfig.Type == RemoteApp, "app type must 'RemoteApp' or 'LocalApp'")
	fmt.Println("LIBGEN CLIENT CONNECTING ...")
	var callable rpcx.Callable = nil
	if gConfig.Type == LocalApp {
		callable = gRpc.NewConnCallable(clientFd, nil)
	} else {
		for {
			var err error = nil
			callable, err = newRemoteCallable(gConfig.Endpoint, time.Second*5)
			if err != nil {
				continue
			}
			//handshake
			out := new(BaseResponse)
			err = callable.Call(ApiCallTimeout, "Handshake", &HandshakeRequest{
				PkgInfo:     gConfig.PkgInfo,
				AccessToken: gConfig.AccessToken,
			}, out)
			if err != nil {
				fmt.Println("LIBGEN CLIENT INIT ERROR , HANDSHAKE FAILED :", err)
				std.CloseIgnoreErr(callable)
				return
			}
			break
		}
	}
	gApiClient.setCallable(callable)
	fmt.Println("LIBGEN CLIENT CONNECTED")
	if gOnConnected != nil {
		gOnConnected()
	}
}

func initWithConfig(config config) {
	initOnce.Do(func() {
		gConfig = config
		doInit()
	})
}

func Cleanup() {
	std.CloseIgnoreErr(GetRawCallable())
	std.CloseIgnoreErr(gRpc)
	gApiClient = nil
}

func doInit() {
	fmt.Println("LIBGEN CLIENT INIT")
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpc failed")
	gRpc = rpc
	gRpc.RegFuncWithName("ControlDevice", onDeviceControl)
	gRpc.RegFuncWithName("DeliveryDeviceStatus", onDeviceStatusDelivery)
	gRpc.RegFuncWithName("Ping", pong)
	gRpc.OnCallableClosed(onCallableClose)
	gRpc.Start()
	gApiClient = NewApiClientImpl()
	fmt.Println("LIBGEN CLIENT INIT SUCCESS")
}

func onCallableClose(callable rpcx.Callable) {
	fmt.Println("LIBGEN RPC DISCONNECTED , RECONNECTING ")
	gApiClient.setCallable(nil)
	go connect()
}

func GetRawCallable() rpcx.Callable {
	std.Assert(gApiClient != nil, "please init first")
	return gApiClient.getCallable()
}

func GetApiClient() RpcApiClient {
	std.Assert(gApiClient != nil, "please init first")
	return gApiClient
}
