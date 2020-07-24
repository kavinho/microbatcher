package microbatcher

import (
	"sync"
)

//Dispatcher is responsible for dispacthing bunch of jobs to batch Processor
//The dispatcching functionality , is put into this struct, to enable multiple dispatchers, in future.
type Dispatcher struct {
	// The long running job is represents by this member.
	ProcessorFn BatchExecuteFn
	// When there is a shutdown signal, this WaitGroup causes code to wait until as dispatches have returned
	flushWait *sync.WaitGroup
}

//Dispatch send long running jgobs to Batch processor
func (dispatcher *Dispatcher) Dispatch(jobsToDispatch []JobWrapper) {

	//This migt need locking...
	if len(jobsToDispatch) == 0 {
		//fmt.Println("dispatcher... No jobs  WaitDone")
		dispatcher.flushWait.Done()
		return
	}

	go func(JobWrappers []JobWrapper) {
		//Have let Waiters know there are pending jobs.
		//dispatcher.flushWait.Add(1)
		//fmt.Printf("dispatcher... Dispatching .. %v \n", len(JobWrappers))
		var jobs []Job = make([]Job, 0)
		//Create a job map to match job result with jJobWrappers
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
			//TODO: if there are JObMap doesnt match result
			//Processor has dropped resutls

			if ok == true {
				//fmt.Printf(" jobW.responseChannel <- r   result is %v \n", r.Result)
				jobW.responseChannel <- r
			}
		}
		//One more dispathc has been done
		//fmt.Println(" Dispatcher.. go..fun Done dispatching")
		dispatcher.flushWait.Done()

	}(jobsToDispatch)

}
