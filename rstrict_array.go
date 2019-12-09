package libgen

import (
	"github.com/gen-iot/std"
	"github.com/pkg/errors"
)

type ArrayLimiter struct {
	baseRestrict
	ElementRestrict Restrict `json:"elementRestrict" validate:"required"`
	FixedLength     int      `json:"fixedLength"`
}

func NewArrayLimiter(name string, elementValidator Restrict) Restrict {
	out := new(ArrayLimiter)
	out.ElementRestrict = elementValidator
	out.FixedLength = -1
	if out.ElementRestrict == nil {
		out.ElementRestrict = NewAnyLimiter(name, false)
	}
	out.baseRestrict = newBaseRestrict(out, name, Array)
	return out.SetDefaultValue(true, std.JsonArray{})
}

func (this *ArrayLimiter) Validate(v interface{}) error {
	if v == nil {
		return errIllegalParams
	}
	va, ok := v.([]interface{})
	if !ok {
		return errIllegalParams
	}
	if this.FixedLength != -1 && len(va) != this.FixedLength {
		return errors.Wrapf(errArrayLengthMismatched, "expect=%d, real=%d", this.FixedLength, len(va))
	}
	for _, it := range va {
		if err := this.ElementRestrict.Validate(it); err != nil {
			return err
		}
	}
	return nil
}
