package read_write

import (
	"bufio"
	"errors"
	o "lecture03_homework/pkg/options"
	"math"
	"os"
)

// minBlockSize is a minimal size of buffer can be created by bufio.NewReaderSize
const minBlockSize = 16

type ReadWrite struct {
	rd        *bufio.Reader
	wr        *bufio.Writer
	readFile  *os.File
	writeFile *os.File
}

func NewReadWrite(opts *o.Options) (rw ReadWrite, err error) {
	rw.rd, rw.readFile, err = rdInit(opts)
	if err != nil {
		return
	}
	rw.wr, rw.writeFile, err = wrInit(opts)
	if err != nil {
		return
	}
	return
}

// Read implements io.Reader
func (rw *ReadWrite) Read(p []byte) (int, error) {
	return rw.rd.Read(p)
}

// Write implements io.Writer
func (rw *ReadWrite) Write(p []byte) (int, error) {
	return rw.wr.Write(p)
}

func (rw *ReadWrite) Flush() error {
	return rw.wr.Flush()
}

func (rw *ReadWrite) Discard(n int) (int, error) {
	return rw.rd.Discard(n)
}

func (rw *ReadWrite) Close() error {
	if err := rw.readFile.Close(); err != nil {
		return err
	}

	if err := rw.writeFile.Close(); err != nil {
		return err
	}
	return nil
}

// rdInit initializing writer according to flags
func rdInit(opts *o.Options) (*bufio.Reader, *os.File, error) {
	var rd *bufio.Reader
	var file *os.File
	if opts.From != "stdin" {
		var err error
		file, err = os.Open(opts.From)
		if err != nil {
			return nil, nil, err
		}
	} else {
		file = os.Stdin
	}

	// creating buffer with size not bigger than opts.BockSize (if possible)
	for {
		rd = bufio.NewReaderSize(file, opts.BlockSize)
		if rd.Size() <= int(math.Max(float64(opts.BlockSize), float64(minBlockSize))) {
			break
		}
	}
	return rd, file, nil
}

// wrInit initializing writer according to flags
func wrInit(opts *o.Options) (*bufio.Writer, *os.File, error) {
	var wr *bufio.Writer
	var File *os.File
	if opts.To != "stdout" {
		if _, err := os.Stat(opts.To); err == nil {
			return nil, nil, errors.New("file " + opts.To + " already exists")
		}
		var err error
		File, err = os.Create(opts.To)
		if err != nil {
			return nil, nil, err
		}
	} else {
		File = os.Stdout
	}
	wr = bufio.NewWriter(File)
	return wr, File, nil
}
