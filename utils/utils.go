package utils

import "fmt"

func MapString[T fmt.Stringer](objects []T) []string {
	var result []string

	for _, obj := range objects {
		result = append(result, obj.String())
	}

	return result
}
