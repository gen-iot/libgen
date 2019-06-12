package libgen

type I32EnumLimiter struct {
	baseRestrict
	Limits []int32 `json:"limits"`
}

func NewI32EnumLimiter(limits []int32) *I32EnumLimiter {
	return &I32EnumLimiter{
		baseRestrict: baseRestrict{
			RestrictType: I32Enum,
		},
		Limits: limits,
	}
}

func (this *I32EnumLimiter) Type() RestrictType {
	return I32Enum
}

func (this *I32EnumLimiter) Validate(v interface{}) error {
	if len(this.Limits) == 0 {
		return errOutOfEnum
	}
	i32, err := any2Int32(v)
	if err != nil {
		return err
	}
	for idx := range this.Limits {
		if this.Limits[idx] == i32 {
			return nil
		}
	}
	return errOutOfEnum
}
