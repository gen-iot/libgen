//+build !server

package libgen

import (
	"errors"
	"fmt"
	"github.com/gen-iot/liblpc/v2"
	"github.com/gen-iot/rpcx/v2"
	"github.com/gen-iot/std"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	initOnce         = sync.Once{}
	gRpcCore         rpcx.Core
	gApiClient       RpcApiClient
	gConfig                 = defaultConfig
	gOnConnected     func() = nil
	ApiCallTimeout          = time.Second * 5
	libgenExitSignal        = make(chan bool)
)

const clientFd = 3

type config struct {
	Type           AppType    `json:"type" validate:"required,oneof=900 901"`
	Endpoint       string     `json:"endpoint"`
	PkgInfo        PkgInfo    `json:"pkgInfo" validate:"-"`
	ApiAccessToken string     `json:"accessToken"`
	LinkMethod     LinkMethod `json:"linkMethod"`
}

var defaultConfig = config{
	Type:           LocalApp,
	Endpoint:       "",
	ApiAccessToken: "",
	PkgInfo: PkgInfo{
		Package: "",
		Name:    "",
	},
}

//noinspection ALL
func InitLocal(onConnected func()) {
	gOnConnected = onConnected
	initWithConfig(defaultConfig)
}

//noinspection ALL
func InitRemote(endPoint string, pkgInfo PkgInfo, apiAccessToken string, onConnected func()) {
	InitRemote2(Handshake, endPoint, pkgInfo, apiAccessToken, onConnected)
}

func InitRemote2(linkMethod LinkMethod, endPoint string, pkgInfo PkgInfo, apiAccessToken string, onConnected func()) {
	gOnConnected = onConnected
	initWithConfig(config{
		LinkMethod:     linkMethod,
		Type:           RemoteApp,
		Endpoint:       endPoint,
		PkgInfo:        pkgInfo,
		ApiAccessToken: apiAccessToken,
	})
}

//noinspection ALL
func Cleanup() {
	close(libgenExitSignal)
	std.CloseIgnoreErr(GetRawCallable())
	std.CloseIgnoreErr(gRpcCore)
	gApiClient = nil
}

//noinspection ALL
func GetRawCallable() rpcx.Callable {
	std.Assert(gApiClient != nil, "please init first")
	return gApiClient.getCallable()
}

//noinspection ALL
func GetApiClient() RpcApiClient {
	std.Assert(gApiClient != nil, "please init first")
	return gApiClient
}

func initWithConfig(config config) {
	initOnce.Do(func() {
		gConfig = config
		doInit()
		go callableWatcher()
	})
}

func doInit() {
	fmt.Printf("LIBGEN CLIENT INIT , MODE=%s\n", AppType2Str(gConfig.Type))
	initSuccessMsg := "LIBGEN CLIENT INIT SUCCESS"
	if gConfig.Type == LocalApp {
		appIdentifier := os.Getenv("X_GEN_APP_IDENTIFIER")
		fmt.Printf("LIBGEN INIT, APP IDENTIFIER=[%s]\n", appIdentifier)
		initSuccessMsg = fmt.Sprintf("%s: APP IDENTIFIER=[%s]", initSuccessMsg, appIdentifier)
	}
	core, err := rpcx.New()
	std.AssertError(err, "new rpc failed")
	gRpcCore = core
	gRpcCore.RegFuncWithName(kCommandDevice, onDeviceCommand)
	gRpcCore.RegFuncWithName(kDeliveryDeviceStatus, onDeviceStatusDelivery)
	gRpcCore.RegFuncWithName(kPing, pong)
	gRpcCore.RegFuncWithName(kInvokeService, onServiceInvoke)
	gRpcCore.RegFuncWithName(kNotifyDeviceIDLE, onDeviceIDLENotify)
	gRpcCore.Start(nil)
	gApiClient = NewApiClientImpl()
	fmt.Println(initSuccessMsg)
}

func waitRemoteConnReady(call *rpcx.SignalCallable, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case err := <-call.ReadySignal():
		//ready!
		return err
	case <-timer.C:
		std.CloseIgnoreErr(call)
		return errors.New("wait for callable ready timeout")
	}
}

func newRemoteCallable(endpoint string) (rpcx.Callable, error) {
	sockAddr, err := liblpc.ResolveTcpAddr(endpoint)
	if err != nil {
		return nil, err
	}
	call, err := rpcx.NewClientStreamCallable(gRpcCore, sockAddr, nil)
	if err != nil {
		return nil, err
	}
	return call, nil
}

func createCallable() (rpcx.Callable, error) {
	std.Assert(strings.Compare(AppType2Str(gConfig.Type), "UNKNOWN") != 0, "unknown app type")
	if gConfig.Type == LocalApp {
		callable := rpcx.NewConnStreamCallable(gRpcCore, clientFd, nil)
		return callable, nil
	} else {
		return newRemoteCallable(gConfig.Endpoint)
	}
}

func doHandshake(call rpcx.Callable) error {
	if gConfig.Type == LocalApp {
		return nil
	}
	linkMethod := gConfig.LinkMethod
	if len(linkMethod) == 0 {
		linkMethod = Handshake
	}
	//handshake
	return call.Call1(ApiCallTimeout, string(linkMethod), &HandshakeRequest{
		PkgInfo:        gConfig.PkgInfo,
		ApiAccessToken: gConfig.ApiAccessToken,
	})
}

func callableWatcher() {
	for {
		//
		select {
		case <-libgenExitSignal:
			log.Println("LIBGEN EXIT...")
			return
		default:
		}
		//
		rcall, err := createCallable()
		std.AssertError(err, "create callable failed")
		call := rpcx.NewSignalCallable(rcall)
		call.Start()
		err = waitRemoteConnReady(call, time.Second*5)
		log.Println("LIBGEN CLIENT CONNECTING ...")
		// connect
		if err != nil {
			log.Println("LIBGEN CLIENT CONNECT FAILED  :", err, " , RECONNECT IN 5s......")
			// connect error , retry after 5s
			time.Sleep(time.Second * 5)
			continue
		}
		// handshake
		if err = doHandshake(call); err != nil {
			log.Println("LIBGEN CLIENT CONNECT FAILED , HANDSHAKE FAILED :", err, " , RECONNECT IN 5s......")
			std.CloseIgnoreErr(call)
			time.Sleep(time.Second * 5)
			continue
		}
		log.Println("LIBGEN CLIENT CONNECTED")
		gApiClient.setCallable(call)
		if gOnConnected != nil {
			gOnConnected()
		}
		// close
		err = <-call.CloseSignal()
		gApiClient.setCallable(nil)
		log.Println("LIBGEN CLIENT DISCONNECTED", " , RECONNECT IN 5s......")
		time.Sleep(time.Second * 5)
	}
}
