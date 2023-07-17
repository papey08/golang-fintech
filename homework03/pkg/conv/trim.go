package conv

import "unicode"

// leftSpacesAmount returns size of space symbols in prefix of the buf
func leftSpacesAmount(buf []byte) int {
	tempStr := string(buf)
	res := 0
	for _, r := range tempStr {
		if unicode.IsSpace(r) {
			res += len(string(r))
		} else {
			break
		}
	}
	return res
}

// rightSpacesAmount returns size of space symbols in suffix of the buf
func rightSpacesAmount(buf []byte) int {
	tempStr := string(buf)
	res := 0
	for _, r := range tempStr {
		if unicode.IsSpace(r) {
			res += len(string(r))
		} else {
			res = 0
		}
	}
	return res
}

func (c *Conv) Trim(buf []byte) []byte {

	if !c.trimmedLeft {
		buf = buf[leftSpacesAmount(buf):]
		if len(buf) != 0 {
			c.trimmedLeft = true
		}
	}

	if !c.trimmedRight {
		if right := rightSpacesAmount(buf); right != len(buf) {
			buf, c.rightSpaces =
				append(c.rightSpaces, buf[:len(buf)-right]...), buf[len(buf)-right:]
		} else {
			c.rightSpaces = append(c.rightSpaces, buf...)
			buf = []byte{}
		}
	}
	return buf
}
