package utils

import (
	"strconv"
	"strings"
)

// StringToSliceUI Convert string to slice of uint
func StringToSliceUI(str string) (res []uint) {
	res = []uint{}
	strSlice := strings.Split(str, ",")
	for _, s := range strSlice {
		if num, err := strconv.Atoi(s); err == nil {
			res = append(res, uint(num))
		}
	}
	return res
}
