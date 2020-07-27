package microbatcher

import (
	"sync"
)

//Dispatcher is responsible for dispacthing bunch of jobs to BatchProcessor
//The dispatching functionality , is put into this struct, to enable multiple dispatchers, in future.
type Dispatcher struct {
	// The long running job is represented by this member.
	ProcessorFn BatchExecuteFn
	// When there is a shutdown signal, this WaitGroup causes code to wait until all dispatches have returned
	flushWait *sync.WaitGroup
}

//Dispatch long running jobs to BatchProcessor
func (dispatcher *Dispatcher) Dispatch(jobsToDispatch []JobWrapper) {

	//TODO: This migt need locking...
	if len(jobsToDispatch) == 0 {
		dispatcher.flushWait.Done()
		return
	}

	go func(JobWrappers []JobWrapper) {
		var jobs []Job = make([]Job, 0)
		//We need a job map to match results to responseChannel
		jobsMap := make(map[string]JobWrapper)

		for _, jobW := range JobWrappers {
			jobs = append(jobs, jobW.theJob)
			jobsMap[jobW.theJob.ID] = jobW
		}
		//Call BatchRrocessor
		results := dispatcher.ProcessorFn(jobs)
		//Let callers know via the channel , that the results are here.
		for _, r := range results {
			jobW, ok := jobsMap[r.JobID]
			//TODO: if there are non matching items in JobMap, Batch Processor has dropped, results
			//Processor has dropped resutls

			if ok == true {
				jobW.responseChannel <- r
			}
		}
		//One more dispatch has been done
		dispatcher.flushWait.Done()

	}(jobsToDispatch)

}
