package libgen

type I32EnumLimiter struct {
	baseRestrict
	Includes []int32 `json:"includes"`
}

func NewI32EnumLimiter(name string, required bool, include ...int32) *I32EnumLimiter {
	return &I32EnumLimiter{
		baseRestrict: baseRestrict{
			RestrictName: name,
			RestrictType: I32Enum,
			Required:     required,
		},
		Includes: include,
	}
}

func NewAnyI32Limiter(name string, required bool) Restrict {
	o := NewI32EnumLimiter(name, required)
	o.RestrictType = I32Any
	return o
}

func (this *I32EnumLimiter) Check(v int32) error {
	includeOk := len(this.Includes) == 0
	for idx := range this.Includes {
		if v == this.Includes[idx] {
			includeOk = true
			break
		}
	}
	if includeOk {
		return nil
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
	return this.Check(v)
}
