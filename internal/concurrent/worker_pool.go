package concurrent

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

type Task = func()
type WorkerPool interface {
	Start() error
	Stop()
	Submit(task Task) error
}

func NewWorkerPool(numberOfWorker, maxQueueLength int) WorkerPool {
	return &workerPoolImpl{
		numberOfWorker: numberOfWorker,
		maxQueueLength: maxQueueLength,
	}
}

type workerPoolImpl struct {
	numberOfWorker int
	maxQueueLength int
	TaskQueue      chan Task
	runningWorkers sync.WaitGroup
}

func (wp *workerPoolImpl) Start() error {
	if wp.TaskQueue != nil {
		return fmt.Errorf("pool has already started")
	}
	if wp.maxQueueLength <= 0 {
		wp.maxQueueLength = math.MaxInt
	}
	if wp.numberOfWorker <= 0 {
		wp.numberOfWorker = runtime.NumCPU()
	}
	wp.TaskQueue = make(chan Task, wp.maxQueueLength)
	for i := 0; i < wp.numberOfWorker; i++ {
		wp.runningWorkers.Add(1)
		go func() {
			defer wp.runningWorkers.Done()
			for task := range wp.TaskQueue {
				if task == nil {
					//task queue has been closed
					return
				}
				task()
			}
		}()
	}
	return nil
}

func (wp *workerPoolImpl) Stop() {
	defer func() {
		// allows multiple routines to call to stop concurrently without panic
		recover()
	}()
	close(wp.TaskQueue)
	wp.TaskQueue = nil
	wp.runningWorkers.Wait()
}

func (wp *workerPoolImpl) Submit(task Task) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic while sending: %v\n", r)
		}
	}()
	select {
	case wp.TaskQueue <- task:
		err = nil
	default:
		err = fmt.Errorf("couldn't submit task to the pool")
	}
	return err
}
