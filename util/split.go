package util

import "strings"

func SplitAndTrim(s, sep string) []string {
	var arr []string
	split := strings.Split(s, sep)

	for _, str := range split {
		s := strings.TrimSpace(str)
		arr = append(arr, s)
	}

	return arr
}
