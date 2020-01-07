package libgen

import (
	"github.com/pkg/errors"
)

type ArrayLimiter struct {
	baseRestrict
	ElementRestrict Restrict `json:"elementRestrict" validate:"required"`
	MinLength       int      `json:"minLength"`
	MaxLength       int      `json:"maxLength"`
	NotEmpty        bool     `json:"notEmpty"`
}

func NewArrayLimiter(name string, elementValidator Restrict) Restrict {
	out := new(ArrayLimiter)
	out.ElementRestrict = elementValidator
	out.MinLength = -1
	out.MaxLength = -1
	out.NotEmpty = false
	if out.ElementRestrict == nil {
		out.ElementRestrict = NewAnyLimiter("validator", false)
	}
	out.baseRestrict = newBaseRestrict(out, name, Array)
	return out.SetDefaultValue(true, []interface{}{})
}

func (this *ArrayLimiter) SetMinLength(minLen int) *ArrayLimiter {
	this.MinLength = minLen
	return this
}

func (this *ArrayLimiter) SetMaxLength(maxLen int) *ArrayLimiter {
	this.MaxLength = maxLen
	return this
}

func (this *ArrayLimiter) Validate(v interface{}) error {
	if v == nil {
		return errIllegalParams
	}
	va, ok := v.([]interface{})
	if !ok {
		return errIllegalParams
	}
	if this.NotEmpty && len(va) == 0 {
		return errArrayEmpty
	}
	if this.MinLength != -1 && len(va) < this.MinLength {
		return errors.Wrapf(errArrayLengthMismatched, "expect >= %d, real=%d", this.MinLength, len(va))
	}
	if this.MaxLength != -1 && len(va) > this.MaxLength {
		return errors.Wrapf(errArrayLengthMismatched, "expect <= %d, real=%d", this.MaxLength, len(va))
	}
	for _, it := range va {
		if err := this.ElementRestrict.Validate(it); err != nil {
			return err
		}
	}
	return nil
}
