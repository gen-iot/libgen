package rpcx

import (
	"gitee.com/Puietel/std"
	"reflect"
	"runtime"
	"strings"
)

func checkInParam(t reflect.Type) {
	inNum := t.NumIn()
	std.Assert(inNum == 1, "func in param len != 1")
	in := t.In(0)
	inKind := in.Kind()
	std.Assert(inKind == reflect.Ptr || inKind == reflect.Struct, "param must be prt of struct")
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

func getValueElement(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}