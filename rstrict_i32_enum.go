package libgen

type I32EnumLimiter struct {
	baseRestrict
	Includes []int32 `json:"includes"`
}

func NewI32EnumLimiter(name string, required bool, include ...int32) *I32EnumLimiter {
	out := new(I32EnumLimiter)
	out.baseRestrict = newBaseRestrict(out, name, I32Enum)
	out.SetRequired(required)
	out.Includes = include
	return out
}

type AnyI32Limiter = I32EnumLimiter

func NewAnyI32Limiter(name string, required bool) *AnyI32Limiter {
	o := NewI32EnumLimiter(name, required)
	o.RestrictType = I32Any
	return o
}

func (this *I32EnumLimiter) Check(v int32) error {
	for idx := range this.Includes {
		if v == this.Includes[idx] {
			return nil
		}
	}
	return errOutOfEnum
}

func (this *I32EnumLimiter) Validate(vx interface{}) error {
	if vx == nil {
		return errIllegalParams
	}
	v, err := any2Int32(vx)
	if err != nil {
		return err
	}
	if this.Type() == I32Any {
		return nil
	}
	return this.Check(v)
}
