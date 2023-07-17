package storage

import (
	"context"
	"errors"
	"sync"
)

const maxWorkersCount = 4

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{maxWorkersCount: maxWorkersCount}
}

type dirStatus int

const (
	inWaiting dirStatus = iota
	inProcess
)

// dirWithStatus is a cover of Dir type where it has status inWaiting or inProcess
type dirWithStatus struct {
	d      Dir
	status dirStatus
}

// AddStatus returns []dirWithStatus based on dirs where every dir has status inWaiting
func addStatus(dirs []Dir) []dirWithStatus {
	dirsWithStatus := make([]dirWithStatus, len(dirs))
	for i, d := range dirs {
		dirsWithStatus[i] = dirWithStatus{
			d:      d,
			status: inWaiting,
		}
	}
	return dirsWithStatus
}

// sendTempRes sends tempRes to tempResChan, returns true if something was sent
func sendTempRes(ctx context.Context, tempResChan chan Result, tempRes Result) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		tempResChan <- tempRes
	}
	return true
}

// sendError sends error to the errChan if need, returns true if something was sent
func sendError(ctx context.Context, errChan chan error, err error) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		if err != nil {
			errChan <- err
			return true
		}
	}
	return false
}

// closeEndChan is for safety endChan closing
func closeEndChan(endChan chan struct{}, isClosed *bool, mu *sync.Mutex) {
	mu.Lock()
	if !*isClosed {
		close(endChan)
		*isClosed = true
	}
	mu.Unlock()
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {

	if a.maxWorkersCount == 0 {
		return Result{}, errors.New("maxWorkersCount is zero")
	}

	// dirs is a slice where new directories will be added in
	dirs := []dirWithStatus{{d, inWaiting}}
	// doneAmount is how many dirs are processed
	doneAmount := 0
	muDirs := new(sync.Mutex)

	// workersCtx is context that signals workers to stop
	workersCtx, cancel := context.WithCancel(ctx)

	// resChan will receive temporary results of calculations from each worker
	resChan := make(chan Result)
	defer close(resChan)

	// errChan will receive errors from each worker
	errChan := make(chan error)
	defer close(errChan)

	// endChan is for sending message from workers to main that calculations are over
	endChan := make(chan struct{})

	// variables for correct closing of endChan
	var isClosed bool
	muIsClosed := new(sync.Mutex)

	for workersCount := 0; workersCount < a.maxWorkersCount; workersCount++ {
		go func() {
			i := 0
			for {
				select {
				case <-workersCtx.Done():
					return
				default:
					// check if all dirs are processed
					muDirs.Lock()
					curDone := doneAmount
					curLen := len(dirs)
					muDirs.Unlock()
					if i >= curLen && curDone == curLen {
						closeEndChan(endChan, &isClosed, muIsClosed)
						return
					} else if i >= curLen {
						continue
					}

					var currentDir dirWithStatus
					muDirs.Lock()
					if dirs[i].status == inWaiting { // getting dir to process
						dirs[i].status = inProcess
						currentDir = dirs[i]
						muDirs.Unlock()
					} else { // go to next dir if current is processing now of already processed
						muDirs.Unlock()
						i++
						continue
					}

					// getting dirs
					tempDirs, tempFiles, lsErr := currentDir.d.Ls(ctx)
					if sendError(ctx, errChan, lsErr) {
						closeEndChan(endChan, &isClosed, muIsClosed)
						return
					}

					// append tempDirs to dirs for processing them in the future
					muDirs.Lock()
					dirs = append(dirs, addStatus(tempDirs)...)
					muDirs.Unlock()

					// calculating total size of all files from currentDir
					var tempFilesSize int64
					for _, f := range tempFiles {
						fSize, statErr := f.Stat(ctx)
						if sendError(ctx, errChan, statErr) {
							closeEndChan(endChan, &isClosed, muIsClosed)
							return
						}

						tempFilesSize += fSize
					}

					// sending result of iteration of goroutine to resChan
					tempRes := Result{
						Size:  tempFilesSize,
						Count: int64(len(tempFiles)),
					}
					sendTempRes(ctx, resChan, tempRes)

					// incrementing counter of processed dirs
					muDirs.Lock()
					doneAmount++
					muDirs.Unlock()
					i++
				}
			}
		}()
	}

	var res Result
	var err error

MainLoop:
	for {
		select {
		case <-ctx.Done():
			cancel()
			closeEndChan(endChan, &isClosed, muIsClosed)
			break MainLoop
		case <-endChan:
			cancel()
			break MainLoop
		case err = <-errChan:
			cancel()
			closeEndChan(endChan, &isClosed, muIsClosed)
			break MainLoop
		case tempRes := <-resChan:
			res.Count += tempRes.Count
			res.Size += tempRes.Size
		}
	}
	if err != nil {
		return Result{}, err
	}
	return res, nil
}
