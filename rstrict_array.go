package libgen

import (
	"github.com/pkg/errors"
	"reflect"
)

type ArrayLimiter struct {
	baseRestrict
	ElementRestrict Restrict `json:"elementRestrict" validate:"required"`
	MinLength       int      `json:"minLength"`
	MaxLength       int      `json:"maxLength"`
	NotEmpty        bool     `json:"notEmpty"`
}

func NewArrayLimiter(name string, elementValidator Restrict) *ArrayLimiter {
	out := new(ArrayLimiter)
	out.ElementRestrict = elementValidator
	out.MinLength = -1
	out.MaxLength = -1
	out.NotEmpty = false
	if out.ElementRestrict == nil {
		out.ElementRestrict = NewAnyLimiter("validator", false)
	}
	out.baseRestrict = newBaseRestrict(out, name, Array)
	out.SetDefaultValue(true, []interface{}{})
	return out
}

func (this *ArrayLimiter) SetMinLength(minLen int) *ArrayLimiter {
	this.MinLength = minLen
	return this
}

func (this *ArrayLimiter) SetMaxLength(maxLen int) *ArrayLimiter {
	this.MaxLength = maxLen
	return this
}

func (this *ArrayLimiter) SetNotEmpty(notEmpty bool) *ArrayLimiter {
	this.NotEmpty = notEmpty
	return this
}

func (this *ArrayLimiter) Validate(v interface{}) error {
	if v == nil {
		return errIllegalParams
	}
	rawValue := reflect.ValueOf(v)
	rawValueKind := rawValue.Kind()
	if rawValueKind != reflect.Array && rawValueKind != reflect.Slice {
		return errIllegalParams
	}

	vaLen := rawValue.Len()
	if this.NotEmpty && vaLen == 0 {
		return errArrayEmpty
	}

	if this.MinLength != -1 && vaLen < this.MinLength {
		return errors.Wrapf(errArrayLengthMismatched, "expect >= %d, real=%d", this.MinLength, vaLen)
	}
	if this.MaxLength != -1 && vaLen > this.MaxLength {
		return errors.Wrapf(errArrayLengthMismatched, "expect <= %d, real=%d", this.MaxLength, vaLen)
	}
	for idx := 0; idx < vaLen; idx++ {
		it := rawValue.Index(idx).Interface()
		if err := this.ElementRestrict.Validate(it); err != nil {
			return err
		}
	}
	return nil
}
