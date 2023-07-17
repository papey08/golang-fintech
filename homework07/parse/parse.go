package parse

import (
	"strconv"
	"strings"
)

type ValidationOperation int

const (
	Wrong ValidationOperation = iota // special value for case when validate tag is invalid
	Length
	In
	Min
	Max
)

// ValidationParams decomposes validate tag to type of operation and args
func ValidationParams(validateTag string) (v ValidationOperation, args any) {
	temp := strings.Split(validateTag, ":")
	if len(temp) <= 1 {
		return Wrong, nil
	}
	switch temp[0] {
	case "len":
		if n, err := strconv.Atoi(temp[1]); err != nil {
			return Wrong, nil
		} else {
			return Length, n
		}
	case "in":
		args = strings.Split(temp[1], ",")
		if args.([]string)[0] == "" {
			return Wrong, nil
		}
		return In, args
	case "min":
		if n, err := strconv.Atoi(temp[1]); err != nil {
			return Wrong, nil
		} else {
			return Min, n
		}
	case "max":
		if n, err := strconv.Atoi(temp[1]); err != nil {
			return Wrong, nil
		} else {
			return Max, n
		}
	default:
		return Wrong, nil
	}
}
