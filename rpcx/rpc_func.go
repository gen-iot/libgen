package rpcx

import (
	"errors"
	"gitee.com/Puietel/std"
	"log"
	"reflect"
)

type rpcFunc struct {
	name      string
	fun       reflect.Value
	inP0Type  reflect.Type
	outP0Type reflect.Type
}

func (this *rpcFunc) decodeInParam(data []byte) (interface{}, error) {
	elementType := this.inP0Type
	isPtr := false
	if this.inP0Type.Kind() == reflect.Ptr {
		elementType = this.inP0Type.Elem()
		isPtr = true
	}
	newOut := reflect.New(elementType).Interface()
	err := std.MsgpackUnmarshal(data, newOut)
	if err != nil {
		return nil, err
	}
	if !isPtr {
		newOut = reflect.ValueOf(newOut).Elem().Interface()
	}
	return newOut, nil
}

func (this *rpcFunc) Call(remoteCallable Callable, inBytes []byte) (outBytes []byte, err error) {
	defer func() {
		panicErr := recover()
		if panicErr == nil {
			return
		}
		log.Println("call error ", panicErr)
		err = errors.New("invoke failed")
	}()
	inParam, err := this.decodeInParam(inBytes)
	if err != nil {
		return nil, err
	}
	paramV := []reflect.Value{reflect.ValueOf(remoteCallable), reflect.ValueOf(inParam)}
	retV := this.fun.Call(paramV)
	if !retV[1].IsNil() {
		err = retV[1].Interface().(error)
	}
	outParam := retV[0].Interface()
	outBytes, err = std.MsgpackMarshal(outParam)
	if err != nil {
		return nil, err
	}
	return outBytes, nil
}
