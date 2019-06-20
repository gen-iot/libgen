package libgen

import (
	"errors"
	"fmt"
	"strings"
)

type StrEnumLimiter struct {
	baseRestrict
	Limits []string `json:"limits" validate:"gte=1"`
}

func NewStrEnumLimiter(limits []string) *StrEnumLimiter {
	return &StrEnumLimiter{
		baseRestrict: baseRestrict{
			RestrictType: StrEnum,
		},
		Limits: limits,
	}
}

func (this *StrEnumLimiter) Validate(v interface{}) error {
	if len(this.Limits) == 0 {
		return errOutOfEnum
	}
	s, err := any2Str(v)
	if err != nil {
		return err
	}
	for idx := range this.Limits {
		if strings.Compare(this.Limits[idx], s) == 0 {
			return nil
		}
	}
	return errOutOfEnum
}

func any2Str(v interface{}) (string, error) {
	switch s := v.(type) {
	case *string:
		return *s, nil
	case string:
		return s, nil
	default:
		return "", errors.New(fmt.Sprintf("%T cant convert to string", v))
	}
}
