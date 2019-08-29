package libgen

import (
	"github.com/gen-iot/std"
	"testing"
)

func must(err error) {
	std.AssertError(err, "failed")
}

func mustErr(err error) {
	std.Assert(err != nil, "failed")
}

func TestI32EnumLimiter(t *testing.T) {
	limiter := NewI32EnumLimiter("test", true, 1, 2, 3, 4, 5)
	must(limiter.Validate(1))
	must(limiter.Validate(2))
	must(limiter.Validate(3))
	mustErr(limiter.Validate(40))
	mustErr(limiter.Validate(50))
}
