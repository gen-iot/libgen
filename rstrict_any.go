package libgen

type AnyLimiter struct {
	baseRestrict
}

func NewAnyLimiter(name string, required bool) *AnyLimiter {
	o := &AnyLimiter{}
	o.RestrictName = name
	o.Required = required
	o.RestrictType = Any
	return o
}

func (this *AnyLimiter) Validate(v interface{}) error {
	return nil
}
