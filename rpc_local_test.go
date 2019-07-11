package libgen

import (
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/libgen/rpcx"
	"gitee.com/SuzhenProjects/libgen/rpcx/middleware"
	"gitee.com/SuzhenProjects/liblpc"
	"log"
	"os"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
)

type testReq struct {
	A  int `validate:"eq=100"`
	B  int
	Tm time.Time
}

type testRsp struct {
	Sum int
}

func sum(ctx rpcx.Context, req *testReq) (*testRsp, error) {
	log.Println("req delta time -> ", time.Now().Sub(req.Tm))
	//headers := ctx.Headers(rpcx.In)
	//if len(headers) != 0 {
	//	log.Println("header -> ", headers)
	//}
	return &testRsp{
		Sum: req.A + req.B,
	}, nil
}

func startLocalRpcService(fd int, wg *sync.WaitGroup) {
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpcx")
	defer std.CloseIgnoreErr(rpc)
	rpc.Start()
	rpc.RegFun(sum)
	rpc.NewConnCallable(fd, nil)
	wg.Wait()
}

func startMockRpcCall(fd int, wg *sync.WaitGroup) {
	defer wg.Done()
	rpc, err := rpcx.New()
	std.AssertError(err, "new rpcx")
	defer std.CloseIgnoreErr(rpc)
	rpc.Start()
	callable := rpc.NewConnCallable(fd, nil,
		middleware.Validate(std.DefaultValidatorEN()),
		middleware.Recover(true),
		middleware.Dump())
	after := time.After(time.Second * 5)
	for {
		select {
		case <-after:
			return
		default:
		}
		out := new(testRsp)
		//header := map[string]string{}
		//for i := 0; i < 5; i++ {
		//	k := fmt.Sprintf("key-%d", i)
		//	v := fmt.Sprintf("val-%d", i)
		//	header[k] = v
		//}
		err = callable.CallWithHeader(time.Second*5, "sum",
			nil,
			&testReq{
				A:  10,
				B:  100,
				Tm: time.Now(),
			}, out)
		std.AssertError(err, "call error")
		std.Assert(out.Sum == 10+100, "result error")
	}
}

func TestRpc(t *testing.T) {
	file, err := os.OpenFile("cpu_prof.prof", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	std.AssertError(err, "create prof failed")
	defer std.CloseIgnoreErr(file)
	err = pprof.StartCPUProfile(file)
	std.AssertError(err, "start profile failed")
	defer pprof.StopCPUProfile()
	//
	fds, err := liblpc.MakeIpcSockpair(true)
	std.AssertError(err, "socketPair error")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startMockRpcCall(fds[1], wg)
	startLocalRpcService(fds[0], wg)
}

func createTraceMiddleware(tag string) rpcx.MiddlewareFunc {
	return func(next rpcx.HandleFunc) rpcx.HandleFunc {
		return func(ctx rpcx.Context) {
			log.Printf("%s In\n", tag)
			req := ctx.Request().(*testReq)
			req.Tm = req.Tm.Add(time.Second * 1)
			next(ctx)
			log.Printf("%s Out\n", tag)
		}
	}
}
