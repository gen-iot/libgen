package libgen

import (
	"errors"
	"reflect"
)

type AnyI32Limiter struct {
	baseRestrict
}

func NewAnyI32Limiter() *AnyLimiter {
	o := &AnyLimiter{}
	o.RestrictType = I32Any
	return o
}

var errCantConvert2I32 = errors.New("param cant convert to int32")

func (this *AnyI32Limiter) Validate(v interface{}) error {
	of := reflect.ValueOf(v)
	if of.Kind() >= reflect.Int && of.Kind() <= reflect.Uint16 {
		return nil
	}
	return errCantConvert2I32
}
