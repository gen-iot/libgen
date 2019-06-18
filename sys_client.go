//+build client

package libgen

import (
	"fmt"
	"gitee.com/Puietel/std"
	"libgen/rpcx"
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
	gRpc.Start()
	gApiClient = new(ApiClientImpl)
	gCallable = gRpc.NewCallable(clientFd, nil)
	fmt.Println("LIBGEN CLIENT INIT SUCCESS")
}

func GetApiClient() *ApiClientImpl {
	return gApiClient
}
