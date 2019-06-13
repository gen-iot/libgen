package libgen

import (
	"fmt"
	"gitee.com/Puietel/std"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"runtime"
	"strings"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type rpcFunc struct {
	name      string
	fun       reflect.Value
	inP0Type  reflect.Type
	outP0Type reflect.Type
}

func (this *rpcFunc) Call(param []interface{}) (out interface{}, err error) {
	defer func() {
		panicErr := recover()
		if panicErr == nil {
			return
		}
		log.Println("call error ", panicErr)
		err = errors.New("invoke failed!")
	}()
	paramV := make([]reflect.Value, 0, len(param))
	for idx := range param {
		paramV = append(paramV, reflect.ValueOf(param[idx]))
	}
	retV := this.fun.Call(paramV)
	fmt.Println(retV)
	if !retV[1].IsNil() {
		err = retV[1].Interface().(error)
	}
	out = retV[0].Interface()
	return

}

type RPC struct {
	rcpFuncMap map[string]*rpcFunc
}

func NewRpc() *RPC {
	return &RPC{
		rcpFuncMap: make(map[string]*rpcFunc),
	}
}

func checkInParam(t reflect.Type) {
	inNum := t.NumIn()
	std.Assert(inNum == 1, "func in param len != 1")
}

func checkOutParam(t reflect.Type) {
	outNum := t.NumOut()
	std.Assert(outNum == 2, "func out param len != 2 ")
	out1 := t.Out(1)
	std.Assert(out1 == typeOfError, "out1 must be error type")
}

func getFuncName(fv reflect.Value) string {
	fname := runtime.FuncForPC(reflect.Indirect(fv).Pointer()).Name()
	idx := strings.LastIndex(fname, ".")
	return fname[idx+1:]
}

func (this *RPC) RegFun(f interface{}) {
	fv, ok := f.(reflect.Value)
	if !ok {
		fv = reflect.ValueOf(f)
	}
	std.Assert(fv.Kind() == reflect.Func, "f not func!")
	fname := getFuncName(fv)
	fvType := fv.Type()
	//check in/out param
	checkInParam(fvType)
	checkOutParam(fvType)
	this.rcpFuncMap[fname] = &rpcFunc{
		name:      fname,
		fun:       fv,
		inP0Type:  fvType.In(0),
		outP0Type: fvType.Out(0),
	}
}

func getValueElement(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

func (this *RPC) MockCall(name string, param interface{}, out interface{}) error {
	outV := reflect.ValueOf(out)
	std.Assert(!outV.IsNil(), "out must not be nil")
	std.Assert(outV.Kind() == reflect.Ptr, "out must be pointer")
	fn, ok := this.rcpFuncMap[name]
	if !ok {
		return errors.New("func not func")
	}
	ret, err := fn.Call([]interface{}{param})
	if err != nil {
		return err
	}
	outV.Elem().Set(getValueElement(reflect.ValueOf(ret)))
	return nil
}
