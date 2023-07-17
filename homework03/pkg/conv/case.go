package conv

import (
	"math"
	"strings"
	"unicode"
	"unicode/utf8"
)

// validSuffix checks if buf ends with correct sequence of bytes
func validSuffix(buf []byte) bool {
	for i := 1; i <= int(math.Min(float64(len(buf)), 4)); i++ {
		if utf8.Valid(buf[len(buf)-i:]) {
			return true
		}
	}
	return false
}

// validPrefix checks if buf starts with correct sequence of bytes
func validPrefix(buf []byte) bool {
	for i := 1; i <= int(math.Min(float64(len(buf)), 4)); i++ {
		if utf8.Valid(buf[:i]) {
			return true
		}
	}
	return false
}

// invalidBoundaries returns if buf is correct utf8 sequence
// and how many bytes are incorrect from the left and from the right side
func invalidBoundaries(buf []byte) (isValid bool, left, right int) {
	for i := 0; i < int(math.Min(float64(len(buf)), 4)); i++ {
		if !validPrefix(buf[i:]) {
			left++
		} else {
			break
		}
	}

	for i := 0; i < int(math.Min(float64(len(buf)), 4)); i++ {
		if !validSuffix(buf[:len(buf)-i]) {
			right++
		} else {
			break
		}
	}

	isValid = left == 0 && right == 0
	return
}

func (c *Conv) Case(buf []byte) []byte {

	// always keeps invalid suffix of buf in invalidBytes
	if len(c.invalidBytes) != 0 {
		isValid, _, _ := invalidBoundaries(buf)
		if !isValid {
			temp := append(c.invalidBytes, buf...)
			_, left, right := invalidBoundaries(temp)
			if right == 0 {
				buf = temp
				c.invalidBytes = []byte{}
			} else if left == 0 {
				buf = temp[:len(temp)-right]
				c.invalidBytes = temp[len(temp)-right:]
			} else {
				c.invalidBytes = temp
				buf = []byte{}
			}
		} else {
			buf = append(c.invalidBytes, buf...)
		}
	} else if _, _, right := invalidBoundaries(buf); right != 0 {
		c.invalidBytes = buf[len(buf)-right:]
		buf = buf[:len(buf)-right]
	}

	switch c.bc {
	case uppercase:
		s := string(buf)
		s = strings.Map(unicode.ToUpper, s)
		return []byte(s)
	case lowercase:
		s := string(buf)
		s = strings.Map(unicode.ToLower, s)
		return []byte(s)
	}
	return buf
}
