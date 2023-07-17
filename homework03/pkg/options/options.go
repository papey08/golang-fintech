package options

import (
	"errors"
	"flag"
	"path/filepath"
	"strings"
)

type Options struct {
	From       string
	To         string
	BlockSize  int
	Offset     int
	Limit      int
	UpperCase  bool
	LowerCase  bool
	TrimSpaces bool
}

const DefaultBufferSize = 4096

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "stdin", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "stdout", "file to write. by default - stdout")

	flag.IntVar(&opts.BlockSize, "block-size", DefaultBufferSize, "size of one chunk, by default 4096")
	flag.IntVar(&opts.Offset, "offset", 0, "amount of bytes to skip when copying")
	flag.IntVar(&opts.Limit, "limit", -1, "amount of bytes to copy")

	var conv string
	flag.StringVar(&conv, "conv", "", "conv arguments separated with commas")

	flag.Parse()

	if opts.From != "stdin" {
		opts.From = filepath.Clean(filepath.Join(opts.From))
	}
	if opts.To != "stdout" {
		opts.To = filepath.Clean(filepath.Join(opts.To))
	}

	errs := make([]string, 0, 3)
	if len(conv) != 0 {
		convArgs := strings.Split(conv, ",")
		for _, c := range convArgs {
			switch c {
			case "upper_case":
				opts.UpperCase = true
			case "lower_case":
				opts.LowerCase = true
			case "trim_spaces":
				opts.TrimSpaces = true
			default:
				errs = append(errs, "unknown conv flag: "+c)
			}
		}
	}
	if opts.LowerCase && opts.UpperCase {
		errs = append(errs, "both upper_case and lower_case")
	}
	if opts.BlockSize == 0 {
		errs = append(errs, "block-size = 0")
	}
	if len(errs) == 0 {
		return &opts, nil
	} else {
		return nil, errors.New(strings.Join(errs, ", "))
	}
}
