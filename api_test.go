package libgen

import (
	"fmt"
	"gitee.com/Puietel/std"
	"log"
	"sync"
	"testing"
	"time"
)

func TestApiClientImpl_Ping(t *testing.T) {
	Init()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		count := 0
		for {
			log.Println("ping test count :", count)
			rsp, err := GetApiClient().Ping(&Ping{Time: time.Now(), Msg: fmt.Sprintf("client ping %d", count)})
			std.AssertError(err, "ping error")
			log.Println("ping res msg >> ", rsp.Msg)
			time.Sleep(time.Millisecond * 50)
			count++
			if count >= 10000 {
				break
			}
		}
		//time.Sleep(time.Second * 60 * 60)
		wg.Done()
		log.Println("ping test over")
	}()
	wg.Wait()
}
