//+build !server

package libgen

import (
	"errors"
	"fmt"
	"github.com/gen-iot/liblpc"
	"github.com/gen-iot/rpcx"
	"github.com/gen-iot/std"
	"log"
	"os"
	"strings"
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

//noinspection ALL
func InitLocal(onConnected func()) {
	initWithConfig(defaultConfig)
	gOnConnected = onConnected
	go callableWatcher()
}

//noinspection ALL
func InitRemote(endPoint string, pkgInfo PkgInfo, accessToken string, onConnected func()) {
	initWithConfig(config{
		Type:        RemoteApp,
		Endpoint:    endPoint,
		PkgInfo:     pkgInfo,
		AccessToken: accessToken,
	})
	gOnConnected = onConnected
	go callableWatcher()
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
	call, err := gRpc.NewClientCallable(sockAddr, nil)
	if err != nil {
		return nil, err
	}
	return call, nil
}

func createCallable() (rpcx.Callable, error) {
	std.Assert(strings.Compare(AppType2Str(gConfig.Type), "UNKNOWN") != 0, "unknown app type")
	if gConfig.Type == LocalApp {
		callable := gRpc.NewConnCallable(clientFd, nil)
		return callable, nil
	} else {
		return newRemoteCallable(gConfig.Endpoint)
	}
}

func doHandshake(call rpcx.Callable) error {
	if gConfig.Type == LocalApp {
		return nil
	}
	//handshake
	return call.Call1(ApiCallTimeout, "Handshake", &HandshakeRequest{
		PkgInfo:     gConfig.PkgInfo,
		AccessToken: gConfig.AccessToken,
	})
}

var libgenExitSignal = make(chan bool)

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

func initWithConfig(config config) {
	initOnce.Do(func() {
		gConfig = config
		doInit()
	})
}

//noinspection ALL
func Cleanup() {
	std.CloseIgnoreErr(GetRawCallable())
	std.CloseIgnoreErr(gRpc)
	gApiClient = nil
}

const (
	// supported func list
	kCommandDevice        = "CommandDevice"
	kDeliveryDeviceStatus = "DeliveryDeviceStatus"
	kPing                 = "Ping"
	kTransportData        = "TransportData"
)

func doInit() {
	fmt.Printf("LIBGEN CLIENT INIT , MODE=%s\n", AppType2Str(gConfig.Type))
	initSuccessMsg := "LIBGEN CLIENT INIT SUCCESS"
	if gConfig.Type == LocalApp {
		appIdentifier := os.Getenv("X_GEN_APP_IDENTIFIER")
		fmt.Printf("LIBGEN INIT, APP IDENTIFIER=[%s]\n", appIdentifier)
		initSuccessMsg = fmt.Sprintf("%s: APP IDENTIFIER=[%s]", initSuccessMsg, appIdentifier)
	}
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpc failed")
	gRpc = rpc

	gRpc.RegFuncWithName(kCommandDevice, onDeviceControl)
	gRpc.RegFuncWithName(kDeliveryDeviceStatus, onDeviceStatusDelivery)
	gRpc.RegFuncWithName(kPing, pong)
	gRpc.RegFuncWithName(kTransportData, onDataTransport)
	gRpc.Start()
	gApiClient = NewApiClientImpl()
	fmt.Println(initSuccessMsg)
}

func GetRawCallable() rpcx.Callable {
	std.Assert(gApiClient != nil, "please init first")
	return gApiClient.getCallable()
}

//noinspection ALL
func GetApiClient() RpcApiClient {
	std.Assert(gApiClient != nil, "please init first")
	return gApiClient
}
