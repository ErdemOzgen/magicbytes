package magicbytes

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"runtime"
	"sync"
)

type Haystack struct {
	filePath  string
	metaslice *[]*Meta
	fn        OnMatchFunc
	wg        *sync.WaitGroup
}

//==========
var FileStatus = map[string]bool{}

//===========
//=====================Custom Error Type===================================================
type MyCustomError struct {
	StatusCode int
	ErrDes     string
}

func (m *MyCustomError) Error() string {
	return fmt.Sprintf(m.ErrDes)
}

//=====================Custom Error Type END=================================================

//======================== search variables check ===========================================
func initalVariableCheck(ctx context.Context, directorypath string, metaslice []*Meta) error {
	if directorypath == "" {
		return &MyCustomError{ErrDes: "empty directory  path", StatusCode: 999}
	}

	if metaslice == nil || len(metaslice) > 1000 {
		return &MyCustomError{ErrDes: "slice limit doesnt meet requirements", StatusCode: 998}
	}

	if ctx == nil {
		return &MyCustomError{ErrDes: "null context values", StatusCode: 997}
	}
	return nil
}

//======================== search variables check ===========================================

//========================================Worker=============================================
//https://gobyexample.com/worker-pools
func worker(ctx context.Context, cancel context.CancelFunc, jobs <-chan *Haystack, wg *sync.WaitGroup) {
	defer wg.Done()
loop:
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			if job == nil {
				continue loop
			}
			meta, err := searchMetasInAFile(ctx, job.filePath, job.metaslice)

			if err != nil {
				log.Println(fmt.Errorf("Meta search has failed path: %s, err: %v", job.filePath, err))

				continue
			}

			if meta != nil && !job.fn(job.filePath, meta.Type) {
				cancel()
			}
			FileStatus[job.filePath] = true
		}

	}
}

//========================================Worker END=============================================
// Search and call on match function

func Search(ctx context.Context, searchDirectory string, metaslice []*Meta, onMatch OnMatchFunc) error {
	err := initalVariableCheck(ctx, searchDirectory, metaslice)
	if err != nil {
		fmt.Printf("Something happen : +v", &MyCustomError{ErrDes: "Something happen", StatusCode: 996})
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan *Haystack, 255)

	var waitGroup sync.WaitGroup
	workerCount := runtime.GOMAXPROCS(0)

	waitGroup.Add(workerCount)
	for i := 1; i <= workerCount; i++ {
		go worker(ctx, cancel, jobs, &waitGroup)
	}

	go func() {
		err := filepath.WalkDir(searchDirectory, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				//Log error and pass
				log.Println(path, err, &MyCustomError{ErrDes: "file walk error countered on path", StatusCode: 955})
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if d == nil || !d.Type().IsRegular() {
				return nil
			}
			FileStatus[path] = false
			jobs <- &Haystack{filePath: path, metaslice: &metaslice, fn: onMatch, wg: &waitGroup}

			return nil
		})
		fmt.Println(len(FileStatus))
		counterMap := map[string]int{}
		for k, _ := range FileStatus {
			if _, ok := counterMap[k]; !ok {
				counterMap[k] = 1
				continue
			}
			counterMap[k]++
			if counterMap[k] > 1 {
				fmt.Println(counterMap[k])
			}
		}

		if err != nil {
			log.Println(err, &MyCustomError{ErrDes: "file walk error countered on path", StatusCode: 954})
		}

		close(jobs)
	}()

	waitGroup.Wait()

	return ctx.Err()
}
