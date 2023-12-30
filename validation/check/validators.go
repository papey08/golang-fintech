package check

func validLen(s string, n int) bool {
	return len(s) == n
}

func validIn[T validatable](val T, vals []T) bool {
	for _, v := range vals {
		if val == v {
			return true
		}
	}
	return false
}

func validMin(n int, min int) bool {
	return n >= min
}

func validMax(n int, max int) bool {
	return n <= max
}

func validLenInterval(s string, a, b int) bool {
	return a <= len(s) && len(s) <= b
}
