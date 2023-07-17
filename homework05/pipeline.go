package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	for _, s := range stages {
		// addition channel for ctx support
		myCh := make(chan any)

		go func(in In, myCh chan any) {
			defer close(myCh)
			// transferring values from in to myCh
			for v := range in {
				select {
				case <-ctx.Done(): // checking if ctx was cancelled
					return
				default:
					myCh <- v
				}
			}
		}(in, myCh)
		in = s(myCh)
	}
	return in
}
