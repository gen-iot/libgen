package libgen

import (
	"gitee.com/Puietel/std"
	"log"
	"testing"
	"time"
)

func TestApiCall(t *testing.T) {
	Init()
	rsp, err := GetApiClient().Ping(&Ping{Time: time.Now(), Msg: "ping"})
	std.AssertError(err, "ping error")
	log.Println("ping res msg >> ", rsp.Msg)
}
