package libgen

import (
	"errors"
	"fmt"
)

func any2Int32(v interface{}) (int32, error) {
	out := int32(0)
	switch num := v.(type) {

	//
	case uint8:
		out = int32(num)
	case uint16:
		out = int32(num)
	case uint:
		out = int32(num)
	case uint32:
		out = int32(num)
	case uint64:
		out = int32(num)
	case *uint8:
		out = int32(*num)
	case *uint16:
		out = int32(*num)
	case *uint:
		out = int32(*num)
	case *uint32:
		out = int32(*num)
	case *uint64:
		out = int32(*num)
	//
	case int8:
		out = int32(num)
	case int16:
		out = int32(num)
	case int:
		out = int32(num)
	case int32:
		out = num
	case int64:
		out = int32(num)
	case float32:
		out = int32(num)
	case float64:
		out = int32(num)
	case *int8:
		out = int32(*num)
	case *int16:
		out = int32(*num)
	case *int:
		out = int32(*num)
	case *int32:
		out = *num
	case *int64:
		out = int32(*num)
	case *float32:
		out = int32(*num)
	case *float64:
		out = int32(*num)
	default:
		return 0, errors.New(fmt.Sprintf("cant convert %T to int32", v))
	}
	return out, nil
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
