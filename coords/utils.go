package main

import (
	"regexp"

	"github.com/nimezhu/bed2x"
)

func RegSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}
func parseRegions(s string) ([]bed2x.Bed3i, bool) {
	arr := RegSplit(s, "[;,]")
	retv := make([]bed2x.Bed3i, len(arr))
	var sign = true
	for i := 0; i < len(arr); i++ {
		a, ok := bed2x.ParseRegion(arr[i])
		if !ok {
			sign = false
		} else {
			retv[i] = a
		}
	}
	return retv, sign
}
