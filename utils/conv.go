package utils

import "strconv"

func GetIntValue(v string, def int) int {
	ret, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return ret
}
