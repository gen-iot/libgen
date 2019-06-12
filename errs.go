package libgen

import "errors"

var ErrIllegalParam = errors.New("illegal param")

var errIllegalParams = ErrIllegalParam
var errValidateNotSupport = errors.New("Restrict.Validate not support")
var errOutOfEnum = errors.New("out of specify enums")
var errOutOfRange = errors.New("out of specify enums")
