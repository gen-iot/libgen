package libgen

import (
	"strings"
)

// 如果include和exclude集合均为空,则该limiter效果等同于AnyStrEnum
type StrEnumLimiter struct {
	baseRestrict
	Includes []string `json:"includes"`
}

func NewStrEnumLimiter(include ...string) *StrEnumLimiter {
	return &StrEnumLimiter{
		baseRestrict: baseRestrict{
			RestrictType: StrEnum,
		},
		Includes: include,
	}
}

func NewAnyStrLimiter() Restrict {
	out := NewStrEnumLimiter()
	out.RestrictType = StrAny
	return out
}

func (this *StrEnumLimiter) Check(v string) error {
	includeOk := len(this.Includes) == 0
	for idx := range this.Includes {
		if strings.Compare(v, this.Includes[idx]) == 0 {
			includeOk = true
			break
		}
	}
	if includeOk {
		return nil
	}
	return errOutOfEnum
}

func (this *StrEnumLimiter) Validate(v interface{}) error {
	if v == nil {
		return errIllegalParams
	}
	s, err := any2Str(v)
	if err != nil {
		return err
	}
	return this.Check(s)
}
