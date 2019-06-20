package libgen

import "errors"

type AnyStrLimiter struct {
	baseRestrict
}

func NewAnyStrLimiter() *AnyLimiter {
	o := &AnyLimiter{}
	o.RestrictType = StrAny
	return o
}

func (this *AnyI32Limiter) AnyStrLimiter(v interface{}) error {
	_, ok := v.(string)
	if !ok {
		_, ok = v.(*string)
	}
	if !ok {
		return errors.New("param not string")
	}
	return nil
}
