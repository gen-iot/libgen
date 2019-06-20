package libgen

type I32RangeLimiter struct {
	baseRestrict
	Gte int32 `json:"gte"`
	Lte int32 `json:"lte" validate:"gtfield=Gte"`
}

func NewI32RangeLimiter(gte, lte int32) *I32RangeLimiter {
	return &I32RangeLimiter{
		baseRestrict: baseRestrict{
			RestrictType: I32Range,
		},
		Gte: gte,
		Lte: lte,
	}
}

func (this *I32RangeLimiter) Validate(v interface{}) error {
	if v == nil {
		return errIllegalParams
	}
	i32, err := any2Int32(v)
	if err != nil {
		return errIllegalParams
	}
	fail := false
	fail = i32 < *this.Gte
	if fail {
		return errOutOfRange
	}
	fail = i32 > *this.Lte
	if fail {
		return errOutOfRange
	}
	return nil

}
