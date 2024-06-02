package bruter

import (
	"context"
	"sync"
	"time"
)

type WorkerPool[T any, R any] interface {
    Start(errorCallbackFn func(input T), successCallbackFn func(result R))
    AddTask(task T)
    Stop()
}

type workerPool[T any, R any] struct {
    maxWorkers int
    workers []readerWorker[T, R]

    // this function will do the actual work
    workerFunc func(T) R
    // this function returns TRUE if reutrn value is indeed an error
    errorCheckFn func(R) bool

    timeout time.Duration
    taskChan chan T
    wg *sync.WaitGroup
}

func NewWorkerPool[T any, R any](maxWorkers int, 
    timeout time.Duration,
    workerFunc func(input T) R,
    errorCheckFn func(result R) bool,
) WorkerPool[T, R] {
    return workerPool[T, R] {
        maxWorkers: maxWorkers,
        workers: make([]readerWorker[T, R], maxWorkers),

        workerFunc: workerFunc,

        errorCheckFn: errorCheckFn,

        timeout: timeout,
        taskChan: make(chan T, maxWorkers),
        wg: &sync.WaitGroup{},
    }
}

func (wp workerPool[T, R])Start(errorCallbackFn func(input T), successCallbackFn func(result R)) {
    for range wp.maxWorkers {
        newWorker := readerWorker[T,R]{}
        newWorker.prepare()

        wp.workers = append(wp.workers, newWorker)
        wp.wg.Add(1)
    }
    context, cancel := context.WithTimeout(context.Background(), wp.timeout)
    defer cancel()
    for _, worker := range wp.workers {
        go worker.run(
            wp.workerFunc,
            wp.taskChan,
            context,
            wp.errorCheckFn, 
            errorCallbackFn,
            successCallbackFn,
            wp.wg,
            )
    }
}

func (wp workerPool[T, R])AddTask(task T) {
    wp.taskChan <- task
}

func (wp workerPool[T, P])Stop() {
    for _, worker := range wp.workers {
        worker.stop()
    }
    wp.wg.Wait()
}
