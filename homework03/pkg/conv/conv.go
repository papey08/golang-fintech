package conv

import (
	o "lecture03_homework/pkg/options"
)

type bufCase int

const (
	defaultcase bufCase = iota
	uppercase
	lowercase
)

type Conv struct {
	bc bufCase // to which case convert the conv

	// field for correct Case method work with symbols of size more than 1 byte
	invalidBytes []byte

	// fields for correct Trim method work with different size of chunks
	trimmedLeft  bool
	trimmedRight bool
	rightSpaces  []byte
}

func NewConv(opts *o.Options) Conv {

	var c Conv

	switch {
	case opts.UpperCase:
		c.bc = uppercase
	case opts.LowerCase:
		c.bc = lowercase
	default:
		c.bc = defaultcase
	}

	c.trimmedLeft = !opts.TrimSpaces
	c.trimmedRight = !opts.TrimSpaces

	return c
}
