package utils

import (
	"fmt"
	"strconv"
)

func MapString[T fmt.Stringer](objects []T) []string {
	var result []string

	for _, obj := range objects {
		result = append(result, obj.String())
	}

	return result
}

func MapInt[T fmt.Stringer](objects []T) []int {
	var result []int

	for _, obj := range objects {
		value, _ := strconv.Atoi(obj.String())
		result = append(result, value)
	}

	return result
}
