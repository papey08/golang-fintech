package main

import (
	"fmt"
	cnv "lecture03_homework/pkg/conv"
	cop "lecture03_homework/pkg/copier"
	o "lecture03_homework/pkg/options"
	rw "lecture03_homework/pkg/read_write"
	"os"
)

func HandleError(err error, s string) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, s, err)
		os.Exit(1)
	}
}

func main() {
	opts, err := o.ParseFlags()
	HandleError(err, "can not parse flags:")

	copier, err := rw.NewReadWrite(opts)
	HandleError(err, "can not create copier:")

	conver := cnv.NewConv(opts)

	err = cop.Copy(&copier, &conver, opts)
	HandleError(err, "can not copy:")
}
