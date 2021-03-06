package libgen

import "github.com/pkg/errors"

var errIllegalParams = errors.New("illegal param")
var errValidateNotSupport = errors.New("Restrict.Validate not support")
var errOutOfEnum = errors.New("out of specify enums")
var errOutOfRange = errors.New("out of specify ranges")
var errValueHasBeenExclude = errors.New("value have been excluded")
var errArrayLengthMismatched = errors.New("array length mismatched")
var errArrayEmpty = errors.New("empty array not allowed here")
