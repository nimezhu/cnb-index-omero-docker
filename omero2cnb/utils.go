package main

import (
	"regexp"
	"strconv"
)

var defaultGenome = "hg38"
var (
	posPattern = regexp.MustCompile("^(\\d+)-(\\d+)$")
)

func parsePos(s string) (int, int, bool) {
	if !posPattern.MatchString(s) {
		return -1, -1, false
	}
	match := posPattern.FindStringSubmatch(s)
	start, err1 := strconv.Atoi(match[1])
	end, err2 := strconv.Atoi(match[2])
	if start > end {
		return -1, -1, false
	}
	if err1 != nil || err2 != nil {
		return -1, -1, false
	}
	return start - 1, end, true
}

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

func parseRegion(s string, name string, parentID string) (*BedURI, bool) {
	a := RegSplit(s, "[:]")
	genome := defaultGenome
	var chr string
	start := 0
	end := 0
	color := "0" //RGB
	ok1 := false
	if len(a) == 2 {
		chr = a[0]
		start, end, ok1 = parsePos(a[1])
	} else if len(a) >= 3 {
		genome = a[0]
		chr = a[1]
		start, end, ok1 = parsePos(a[2])
	} else if len(a) >= 4 {
		color = a[3]
	}
	if !ok1 {
		return nil, false
	}
	return &BedURI{genome, chr, start, end, name, color, parentID}, true
}

func parseRegions(s string, prefix string, parentID string) ([]*BedURI, bool) {
	arr := RegSplit(s, "[;,]")
	retv := make([]*BedURI, len(arr))
	var sign = true
	var l = len(arr)
	for i := 0; i < l; i++ {
		var name string
		if l > 1 {
			name = prefix + "_" + strconv.Itoa(i)
		} else {
			name = prefix
		}
		a, ok := parseRegion(arr[i], name, parentID)
		if !ok {
			retv[i] = nil
			sign = false
		} else {
			retv[i] = a
		}
	}
	return retv, sign
}
