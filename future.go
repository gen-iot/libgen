package libgen

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var ErrFutureWaitTimeout = errors.New("future timeout")
var ErrFutureAlreadyInWait = errors.New("future already waiting")
var ErrPromiseAlreadySatisfied = errors.New("promise already satisfied")
var ErrFutureAlreadyRetrieved = errors.New("future already retrieved")

type Future interface {
	Wait(d time.Duration) (interface{}, error)
}

type Promise interface {
	SetResult(interface{}, error) error
	GetFuture() (Future, error)
}

type futureImpl struct {
	promise *promiseImpl
	stat    int32
}

func newFutureImpl(p *promiseImpl) *futureImpl {
	return &futureImpl{
		promise: p,
		stat:    0,
	}
}

func (this *futureImpl) Wait(d time.Duration) (interface{}, error) {
	if !atomic.CompareAndSwapInt32(&this.stat, 0, 1) {
		return nil, ErrFutureAlreadyInWait
	}
	timeoutCtx, cancelFn := context.WithTimeout(context.Background(), d)
	defer cancelFn()
	select {
	case <-timeoutCtx.Done():
		return nil, ErrFutureWaitTimeout //reach timeout
	case <-this.promise.done:
		promise := this.promise
		return promise.resultData, promise.resultErr
	}
	panic("unreachable code")
}

type promiseImpl struct {
	futureRetriedFlag int32
	resultSetFlag     int32
	done              chan bool
	resultData        interface{}
	resultErr         error
}

func NewPromise() Promise {
	return &promiseImpl{
		futureRetriedFlag: 0,
		resultSetFlag:     0,
		done:              make(chan bool),
		resultData:        nil,
		resultErr:         nil,
	}
}

func (this *promiseImpl) SetResult(result interface{}, resultErr error) error {
	if !atomic.CompareAndSwapInt32(&this.resultSetFlag, 0, 1) {
		return ErrPromiseAlreadySatisfied
	}
	this.resultData = result
	this.resultErr = resultErr
	close(this.done)
	return nil
}

func (this *promiseImpl) GetFuture() (Future, error) {
	if !atomic.CompareAndSwapInt32(&this.futureRetriedFlag, 0, 1) {
		return nil, ErrFutureAlreadyRetrieved
	}
	return newFutureImpl(this), nil
}
