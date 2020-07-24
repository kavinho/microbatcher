package microbatcher

import (
	"fmt"
	"sync"
	"time"
)

type sampleProcessor struct {
	Mutex *sync.RWMutex
	Delay time.Duration
}

func (sp *sampleProcessor) Execute(jobs []Job) []JobResult {
	sp.Mutex.Lock()
	defer sp.Mutex.Unlock()
	results := make([]JobResult, len(jobs))
	//simulating a long running process
	time.Sleep(sp.Delay)
	fmt.Printf("Processor.execute() Received Jobs... %v \n ", len(jobs))
	for _, job := range jobs {
		results = append(results, JobResult{Result: job.Param * 2, JobID: job.ID})
	}
	return results
}

func NewSampleProceesor(delay time.Duration) *sampleProcessor {

	return &sampleProcessor{
		Mutex: &sync.RWMutex{},
		Delay: delay,
	}

}
