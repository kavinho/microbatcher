package microbatcher

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDispatcherOneJob(t *testing.T) {
	processor := NewSampleProceesor(time.Millisecond * 300)
	waitSignal := &sync.WaitGroup{}
	dispatcher := &Dispatcher{
		ProcessorFn: processor.Execute,
		flushWait:   waitSignal,
	}
	job := Job{Param: 1, ID: "99"}
	channel := make(chan JobResult)
	jobw := JobWrapper{theJob: job, responseChannel: channel}

	jobsw := []JobWrapper{jobw}
	waitSignal.Add(1)
	dispatcher.Dispatch(jobsw)
	jresult := <-channel

	assert.Equal(t, jresult.Result, 2, "Expected 2")

}
func TestDispatcherTwoJobs(t *testing.T) {
	fmt.Printf("Processor.execute() Received Jobs... %v \n ", 99)
	processor := NewSampleProceesor(time.Millisecond * 300)

	waitSignal := &sync.WaitGroup{}
	dispatcher := &Dispatcher{
		ProcessorFn: processor.Execute,
		flushWait:   waitSignal,
	}
	job1 := Job{Param: 1, ID: "99"}
	job2 := Job{Param: 2, ID: "100"}
	channel := make(chan JobResult)
	jobw1 := JobWrapper{theJob: job1, responseChannel: channel}
	jobw2 := JobWrapper{theJob: job2, responseChannel: channel}

	jobsw := []JobWrapper{jobw1, jobw2}

	waitSignal.Add(1)
	dispatcher.Dispatch(jobsw)
	jresult1 := <-channel
	jresult2 := <-channel

	assert.Equal(t, jresult1.Result, 2, "Expected 2")
	assert.Equal(t, jresult2.Result, 4, "Expected 4")

}

func TestDispatcherGoroutineJobs(t *testing.T) {
	fmt.Printf("Processor.execute() Received Jobs... %v \n ", 99)
	processor := NewSampleProceesor(time.Millisecond * 300)
	waitSignal := &sync.WaitGroup{}
	dispatcher := &Dispatcher{
		ProcessorFn: processor.Execute,
		flushWait:   waitSignal,
	}
	job1 := Job{Param: 1, ID: "99"}
	job2 := Job{Param: 2, ID: "100"}
	channel := make(chan JobResult)
	jobw1 := JobWrapper{theJob: job1, responseChannel: channel}
	jobw2 := JobWrapper{theJob: job2, responseChannel: channel}

	jobsw := []JobWrapper{jobw1, jobw2}

	for i := range []int{0, 1} {
		go func(jbw JobWrapper) {
			waitSignal.Add(1)
			dispatcher.Dispatch([]JobWrapper{jbw})
		}(jobsw[i])

	}

	jresult1 := <-channel
	jresult2 := <-channel

	assert.True(t, jresult1.Result == 2 || jresult1.Result == 4, "Expected result not detected")
	assert.True(t, jresult2.Result == 2 || jresult2.Result == 4, "Expected result not detected")
	assert.NotEqual(t, jresult1.Result, jresult2.Result, "Dipatcher should not return mixed up results")

}
