package rpcx

import (
	"errors"
	"log"
	"reflect"
)

type rpcFunc struct {
	name      string
	fun       reflect.Value
	inP0Type  reflect.Type
	outP0Type reflect.Type
}

func (this *rpcFunc) buildInvoke(ctx *contextImpl) HandleFunc {
	dir := In
	var nextFn func(ctx *contextImpl) = nil
	return func(_ Context) error {
		switch dir {
		case In:
			nextFn = this.call(ctx)
			dir = Out
		case Out:
			if nextFn != nil {
				nextFn(ctx)
			}
		}
		return nil
	}
}

func (this *rpcFunc) decodeInParam(data []byte) (interface{}, error) {
	elementType := this.inP0Type
	isPtr := false
	if this.inP0Type.Kind() == reflect.Ptr {
		elementType = this.inP0Type.Elem()
		isPtr = true
	}
	newOut := reflect.New(elementType).Interface()
	err := gRpcSerialization.UnMarshal(data, newOut)
	if err != nil {
		return nil, err
	}
	if !isPtr {
		newOut = reflect.ValueOf(newOut).Elem().Interface()
	}
	// todo replace with validate middleware
	//err = std.ValidateStruct(newOut)
	//if err != nil {
	//	return nil, err
	//}
	return newOut, nil
}

var errInvokeErr = errors.New("invoke failed")

func (this *rpcFunc) call(ctx *contextImpl) (outFn func(ctx *contextImpl)) {
	inParam, err := this.decodeInParam(ctx.inMsg.Data)
	if err != nil {
		ctx.SetError(err)
		return nil
	}
	ctx.SetRequest(inParam)
	ctx.setDirection(In)
	return func(ctx *contextImpl) {
		defer func() {
			//
			panicErr := recover()
			if panicErr == nil {
				return
			}
			outFn = nil
			log.Println("call error ", panicErr)
			ctx.SetError(errInvokeErr)
		}()
		inParam = ctx.Request()
		paramV := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(inParam)}
		retV := this.fun.Call(paramV)
		if !retV[1].IsNil() {
			err = retV[1].Interface().(error)
			ctx.SetError(err)
			return
		}
		outParam := retV[0].Interface()
		ctx.SetResponse(outParam)
		ctx.setDirection(Out)
	}
}
