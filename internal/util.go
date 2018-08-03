package internal

import (
	"strconv"
	"strings"
)

func StringToInt64Slice(src string) ([]int64, error) {
	var value []int64

	if src == "{}" {
		return value, nil
	}

	src = strings.Trim(src, "{}")
	for _, s := range strings.Split(src, ",") {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		value = append(value, i)
	}
	return value, nil
}
