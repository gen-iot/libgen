//+build !server

package libgen

import (
	"fmt"
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/libgen/rpcx"
	"gitee.com/SuzhenProjects/liblpc"
	"github.com/pkg/errors"
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
	connectToGen()
	fmt.Println("LIBGEN CLIENT INIT SUCCESS")
}
func connectToGen() {
	var callable rpcx.Callable = nil
	var err error = nil
	i := 0
	for {
		callable, err = newCallable(gConfig)
		if err != nil {
			fmt.Println("LIBGEN CLIENT INIT ERROR , CONNECT ERROR :", err)
			i++
			if i >= 100 {
				panic("LIBGEN CLIENT INIT ERROR , AFTER 100 TIMES CONNECT RETRY")
			}
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

var errUnknownAppType = errors.New("unknown app type , app type must 'RemoteApp' or 'LocalApp'")

func newCallable(conf config) (rpcx.Callable, error) {
	if conf.Type == RemoteApp {
		sockFd, err := liblpc.NewConnFd(conf.Endpoint)
		if err != nil {
			return nil, err
		}
		cancelFn, future := gRpc.NewClientCallable(int(sockFd), nil)
		// sync wait
		data, err := future.WaitData(time.Second * 5)
		if err != nil {
			cancelFn()
			return nil, err
		}
		return data.(rpcx.Callable), nil

	} else if conf.Type == LocalApp {
		callable := gRpc.NewConnCallable(clientFd, nil)
		return callable, nil
	}
	return nil, errUnknownAppType
}

func onCallableClose(callable rpcx.Callable) {
	fmt.Println("LIBGEN RPC DISCONNECTED ,RECONNECTING")
	//todo reconnect
}

func GetRawCallable() rpcx.Callable {
	std.Assert(gApiClient != nil, "please init first")
	return gApiClient.getCallable()
}

func GetApiClient() RpcApiClient {
	std.Assert(gApiClient != nil, "please init first")
	return gApiClient
}
