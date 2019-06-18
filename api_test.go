package libgen

import (
	"gitee.com/Puietel/std"
	"libgen/rpcx"
	"testing"
	"time"
)

func TestApiCall(t *testing.T) {
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpcx")
	defer std.CloseIgnoreErr(rpc)
	rpc.Start()
	callable := rpc.NewCallable(3, nil)
	err = callable.Call(time.Second*5, "", "", "")
	std.AssertError(err, "call ")
}
