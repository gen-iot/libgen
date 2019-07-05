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
var gRwLock = &sync.RWMutex{}
var gConfig = defaultConfig

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
	std.CloseIgnoreErr(GetRawCallable())
	std.CloseIgnoreErr(gRpc)
	gApiClient = nil
}

func doInit(config config) {
	fmt.Println("LIBGEN CLIENT INIT")
	gConfig = config
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
	gApiClient = NewApiClientImpl()
	go connectToGen()
	fmt.Println("LIBGEN CLIENT INIT SUCCESS")
}
func connectToGen() {
	var callable rpcx.Callable = nil
	var err error = nil
	for {
		callable, err = newCallable(gConfig)
		if err != nil {
			fmt.Println("LIBGEN CLIENT INIT ERROR , CONNECT ERROR :", err)
			time.Sleep(time.Second * 5)
			continue
		}
		break
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
	gApiClient.setCallable(callable)
}

func newCallable(conf config) (callable rpcx.Callable, err error) {
	if conf.Type == RemoteApp {
		sockFd, err := liblpc.NewConnFd(conf.Endpoint)
		if err != nil {
			return
		}
		//std.AssertError(err, "connect err")
		callable = gRpc.NewCallable(int(sockFd), nil)
		return
	} else {
		callable = gRpc.NewCallable(clientFd, nil)
	}
	return
}

func onCallableClose(callable rpcx.Callable) {
	fmt.Println("LIBGEN RPC DISCONNECTED ,RECONNECTING")
	go connectToGen()
}

func GetRawCallable() rpcx.Callable {
	std.Assert(gApiClient != nil, "please init first")
	return gApiClient.getCallable()
}

func GetApiClient() RpcApiClient {
	std.Assert(gApiClient != nil, "please init first")
	return gApiClient
}
