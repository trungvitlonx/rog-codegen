package util

import (
	"fmt"
)

func ToSliceOfString[T interface{}](s []T) []string {
	r := make([]string, 0, len(s))
	for _, v := range s {
		r = append(r, fmt.Sprintf("%v", v))
	}

	return r
}

func SliceToMap(items []string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, item := range items {
		m[item] = true
	}

	return m
}

func SliceContains(items []string, value string) bool {
	m := SliceToMap(items)
	return m[value]
}
