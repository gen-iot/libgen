package middleware

import (
	"gitee.com/Puietel/std"
	"gitee.com/SuzhenProjects/libgen/rpcx"
	"log"
)

func Validate(v std.Validator) rpcx.MiddlewareFunc {
	std.Assert(v != nil, "validator is nil")
	return func(next rpcx.HandleFunc) rpcx.HandleFunc {
		return func(ctx rpcx.Context) {
			err := ctx.Error()
			if err != nil {
				log.Println("validate req abort ,caused by:", err)
				return
			}
			err = v.Validate(ctx.Request())
			if err != nil {
				log.Println("validate req err")
				ctx.SetError(err)
				return
			}
			//
			next(ctx)
			//
			if ctx.Error() != nil {
				log.Println("validate rsp abort ,caused by:", err)
				return
			}
			err = v.Validate(ctx.Response())
			if err != nil {
				log.Println("validate rsp err")
				ctx.SetError(err)
				return
			}
		}
	}
}
