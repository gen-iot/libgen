package rpcx

import (
	"gitee.com/Puietel/std"
	"reflect"
	"runtime"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfCallable = reflect.TypeOf((*Callable)(nil)).Elem()

func checkInParam(t reflect.Type) {
	inNum := t.NumIn()
	std.Assert(inNum == 2, "func in1 param len != 1")
	in0 := t.In(0)
	std.Assert(in0 == typeOfCallable, "param[0] must be callable")
	in1 := t.In(1)
	in1Kind := in1.Kind()
	std.Assert(in1Kind == reflect.Ptr || in1Kind == reflect.Struct, "param[1] must be prt of struct")
}

func checkOutParam(t reflect.Type) {
	outNum := t.NumOut()
	std.Assert(outNum == 2, "func musts have two out_param")
	out1 := t.Out(1)
	std.Assert(out1 == typeOfError, "out_param[1].type must be `error`")
}

func getFuncName(fv reflect.Value) string {
	//fname := runtime.FuncForPC(reflect.Indirect(fv).Pointer()).Name()
	//idx := strings.LastIndex(fname, ".")
	//return fname[idx+1:]
	return runtime.FuncForPC(fv.Pointer()).Name()
}

func getValueElement(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}
