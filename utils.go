package libgen

import (
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
