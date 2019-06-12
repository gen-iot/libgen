package libgen


type AnyLimiter struct {
	baseRestrict
}

func NewAnyLimiter() *AnyLimiter {
	o := &AnyLimiter{}
	o.RestrictType = Any
	return o
}

func (this *AnyLimiter) Validate(v interface{}) error {
	return nil
}
