package utils

import (
	"strings"
)

func ConvertStringToArray(input string) []string {
	input = strings.Trim(input, `"`)

	result := strings.Split(input, `","`)

	return result
}
