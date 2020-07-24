package microbatcher

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMicroShutDown(t *testing.T) {
	processor := NewSampleProceesor(time.Millisecond * 300)
	fmt.Println("start run  ")
	micro := NewMicroBatcher(processor.Execute, 2, time.Millisecond*100)
	micro.Shutdown()
	//No no panic is a pass
}

func TestMicroOneJobBatchSizeTriggers(t *testing.T) {
	fmt.Println("start run ")
	processor := NewSampleProceesor(time.Millisecond * 300)
	//set to a big value so it never ticks
	cycle := time.Second * 500
	micro := NewMicroBatcher(processor.Execute, 1, cycle)
	//var wg sync.WaitGroup
	//wg.Add(1)
	jb := Job{32, "AAAA"}
	result, _ := micro.Run(jb)
	//wg.Done()
	micro.Shutdown()
	assert.Equal(t, 64, result.Result, "Expected different result")
	fmt.Println("request to shutdown ")
}

func TestMicroOneJobTimerTriggers(t *testing.T) {
	fmt.Println("start run ")
	processor := NewSampleProceesor(time.Millisecond * 300)
	//set to a big value so it never ticks
	cycle := time.Microsecond * 100
	micro := NewMicroBatcher(processor.Execute, 500, cycle)
	//var wg sync.WaitGroup
	//wg.Add(1)
	jb := Job{32, "AAAA"}
	result, _ := micro.Run(jb)
	fmt.Println("Request to shutDown ")
	micro.Shutdown()
	assert.Equal(t, 64, result.Result, "Expected different result")
	fmt.Println("request to shutdown ")
}

func TestMicroOneJobFasterMetronome(t *testing.T) {
	fmt.Println("start run ")
	processor := NewSampleProceesor(time.Millisecond * 300)
	micro := NewMicroBatcher(processor.Execute, 2, time.Millisecond*150)
	var wg sync.WaitGroup
	wg.Add(1)
	jb := Job{32, "AAAAA"}
	fmt.Println("calling run ")
	result, _ := micro.Run(jb)
	wg.Done()
	fmt.Println("called run ")
	assert.Equal(t, 64, result.Result, "Expected different result")
	fmt.Println("request to shutdown ")
	micro.Shutdown()
	fmt.Println("ON Wait ")
	wg.Wait()
}

func TestMicroOneJobFasterProcessor(t *testing.T) {
	fmt.Println("start run ")
	processor := NewSampleProceesor(time.Millisecond * 300)
	micro := NewMicroBatcher(processor.Execute, 2, time.Millisecond*250)
	var wg sync.WaitGroup
	wg.Add(1)
	jb := Job{32, "AAAAA"}
	fmt.Println("calling run ")
	result, _ := micro.Run(jb)
	wg.Done()
	fmt.Println("called run ")
	assert.Equal(t, 64, result.Result, "Expected different result")
	fmt.Println("request to shutdown ")
	micro.Shutdown()
	fmt.Println("ON Wait ")
	wg.Wait()

}

func TestMicroMultiGorotineJobs(t *testing.T) {
	mutex := sync.RWMutex{}

	processor := NewSampleProceesor(time.Millisecond * 30)
	var wg sync.WaitGroup
	insertJobs := 3005
	results := make([]*JobResult, 0)
	micro := NewMicroBatcher(processor.Execute, 20, time.Millisecond*100)
	wg.Add(insertJobs)
	for i := 1; i <= insertJobs; i++ {
		go func(number int) {
			jb := Job{Param: number, ID: fmt.Sprintf("TMMGJ-%d", number)}
			result, _ := micro.Run(jb)
			mutex.Lock()
			results = append(results, result)
			mutex.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	micro.Shutdown()
	assert.Equal(t, insertJobs, len(results), "Not all resutls came back.")
}
