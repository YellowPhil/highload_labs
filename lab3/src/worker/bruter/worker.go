package bruter

import (
	"context"
	"log"
	"sync"
)

type worker struct {
    quitChan chan bool
}
func (w worker) prepare() {
    w.quitChan = make(chan bool, 1)
}

func (w worker) stop() {
    w.quitChan <- true
    close(w.quitChan)
}

// this type of worker will do the actual job: call the worker func provided and check the results, 
type readerWorker[T any, R any] struct {
    worker
}

func (rw *readerWorker[T, R])run(
    runner func(T) R,
    inputChan chan T,
    ctx context.Context,
    errorCheckFn func(result R) bool,
    errorCallbackFn(func (input T)), 
    successCallbackFn(func (r R)),
    wg *sync.WaitGroup,
) {
    for {
        select {
        case input := <- inputChan:
            res := runner(input)
            if ctx.Err() == context.DeadlineExceeded {
                log.Println("Running job takes longer than should. Exitting now...")
                errorCallbackFn(input)
                break
            }
            if err := errorCheckFn(res); err == true {
                log.Println("Reader worker received an error in response")
                errorCallbackFn(input)
            } else {
                successCallbackFn(res)
            }
        case <- rw.quitChan:
        log.Println("Shutting down reader worker")
        wg.Done()
        return
    }
    }
}
