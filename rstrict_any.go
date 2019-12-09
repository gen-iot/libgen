package libgen

type AnyLimiter struct {
	baseRestrict
}

func NewAnyLimiter(name string, required bool) *AnyLimiter {
	o := &AnyLimiter{}
	o.baseRestrict = newBaseRestrict(o, name, Any)
	return o
}

func (this *AnyLimiter) Validate(v interface{}) error {
	return nil
}
