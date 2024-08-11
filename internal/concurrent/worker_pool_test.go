package concurrent

import (
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func Test_workerPoolImpl(t *testing.T) {
	const (
		NUMBER_OF_WORKERS = 2
		MAX_QUEUE_SIZE    = 5
	)
	wp := NewWorkerPool(2, 5)
	completedJobCount := atomic.Int32{}
	cantSubmittedJobCount := 0

	err := wp.Start()
	assert.NoError(t, err, "worker pool must start normally")

	err = wp.Start()
	assert.Error(t, err, "worker pool already running then cannot start again")

	for i := 0; i < MAX_QUEUE_SIZE+3; i++ {
		err := wp.Submit(func() {
			completedJobCount.Add(1)
			time.Sleep(10 * time.Millisecond)
		})
		if err != nil {
			cantSubmittedJobCount++
		}
	}

	time.Sleep(10 * time.Millisecond)

	// for scheduling time slice may not be equal between different go routine,
	// so we cannot make sure 100% the cantSubmittedJobCount always equals to 3. but 99.99999% it will be 3
	assert.Greater(t, cantSubmittedJobCount, 0, "Number of failed to submit job must be greater than 0")
	assert.GreaterOrEqual(t, completedJobCount.Load(), int32(NUMBER_OF_WORKERS),
		"Number of completed jobs must be greater or equal to the number of workers")

	longTaskCompleted := false
	wp.Submit(func() {
		time.Sleep(50 * time.Millisecond)
		longTaskCompleted = true
	})
	time.Sleep(5 * time.Millisecond)
	wp.Stop()
	assert.True(t, longTaskCompleted, "if the task is already running, it must be waited until its done")
	defer func() {
		r := recover()
		if r != nil {
			assert.Fail(t, "multiple go routines can call to stop concurrently,"+
				" then shouldn't get panic in case of stopping multiple times")
		}
	}()
	wp.Stop()
}
