package copier

import (
	"errors"
	"io"
	o "lecture03_homework/pkg/options"
	"strconv"
)

type Copier interface {
	io.ReadWriter
	Discard(n int) (int, error)
	Flush() error
	Close() error
}

type Convert interface {
	Case(buf []byte) []byte
	Trim(buf []byte) []byte
}

func Copy(cop Copier, cnv Convert, opts *o.Options) error {
	readedBytes := 0
	limitReached := false

	// skipping first opts.Offset bytes of the input
	_, err := cop.Discard(opts.Offset)
	if err != nil {
		return err
	}

	for {
		buf := make([]byte, opts.BlockSize)
		n, err := cop.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		} else if n == 0 {
			continue
		}
		buf = buf[:n]

		// applying upperCase or lowerCase to buf
		buf = cnv.Case(buf)

		// checking if opts.Limit is reached
		readedBytes += len(buf)
		if opts.Limit != -1 && readedBytes > opts.Limit {
			buf = buf[:len(buf)-(readedBytes-opts.Limit)]
			limitReached = true
		}

		buf = cnv.Trim(buf)

		n, err = cop.Write(buf)
		if err != nil {
			return err
		} else if n < len(buf) {
			return errors.New("write error: " + strconv.Itoa(len(buf)-n) + " bytes lost")
		}
		if limitReached {
			break
		}
	}
	if err = cop.Flush(); err != nil {
		return err
	}
	return cop.Close()
}
