package libgen

import (
	"strings"
)

// 如果include和exclude集合均为空,则该limiter效果等同于AnyStrEnum
type StrEnumLimiter struct {
	baseRestrict
	Includes []string `json:"includes"`
}

func NewStrEnumLimiter(name string, required bool, include ...string) Restrict {
	out := &StrEnumLimiter{
		Includes: include,
	}
	out.baseRestrict = newBaseRestrict(out, name, StrEnum)
	out.SetRequired(required)
	return out
}

type AnyStrLimiter = StrEnumLimiter

func NewAnyStrLimiter(name string, required bool) Restrict {
	out := NewStrEnumLimiter(name, required).(*StrEnumLimiter)
	out.RestrictType = StrAny
	return out
}

func (this *StrEnumLimiter) Check(v string) error {
	for idx := range this.Includes {
		if strings.Compare(v, this.Includes[idx]) == 0 {
			return nil
		}
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
	if this.RestrictType == StrAny {
		return nil
	}
	return this.Check(s)
}
