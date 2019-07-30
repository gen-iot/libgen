package libgen

import (
	"errors"
	"fmt"
	"strconv"
)

func any2Int32(v interface{}) (int32, error) {
	i, err := strconv.Atoi(fmt.Sprintf("%v", v))
	if err != nil {
		return 0, err
	}
	return int32(i), nil
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
