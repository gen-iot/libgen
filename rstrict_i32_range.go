package libgen

type I32RangeLimiter struct {
	baseRestrict
	Gte      int32   `json:"gte"`
	Lte      int32   `json:"lte" validate:"gtfield=Gte"`
	Addition []int32 `json:"addition"` // 允许出现在[gte,lte]之外的数值
	Exclude  []int32 `json:"exclude"`  // 不允许出现的值,即使在[gte,lte]范围之内
}

func NewI32RangeLimiter(name string, required bool, gte, lte int32, addition []int32, exclude []int32) *I32RangeLimiter {
	out := &I32RangeLimiter{
		Gte:      gte,
		Lte:      lte,
		Addition: addition,
		Exclude:  exclude,
	}
	out.baseRestrict = newBaseRestrict(out, name, I32Range)
	out.SetRequired(required)

	return out
}

func (this *I32RangeLimiter) Validate(v interface{}) error {
	if v == nil {
		return errIllegalParams
	}
	i32, err := any2Int32(v)
	if err != nil {
		return errIllegalParams
	}
	//
	// check extra
	for _, i := range this.Addition {
		if i32 == i {
			return nil
		}
	}
	//
	// check exclude
	for _, i := range this.Exclude {
		if i32 == i {
			return errValueHasBeenExclude
		}
	}
	//
	if i32 < this.Gte || i32 > this.Lte {
		return errOutOfRange
	}
	return nil
}
