package libgen

import "errors"

var errIllegalParams = errors.New("illegal param")
var errValidateNotSupport = errors.New("Restrict.Validate not support")
var errOutOfEnum = errors.New("out of specify enums")
var errOutOfRange = errors.New("out of specify enums")
var errValueHasBeenExclude = errors.New("value have been excluded")
