package microbatcher

import (
	"sync"
	"time"
)

//ErrShutDown an Internal error indicator
type ErrShutDown int

func (e ErrShutDown) Error() string {
	return "Batcher has lready shutdown."
}

//MicroBatcher  an implementation of a micro batcher in GO
type MicroBatcher struct {
	// how many items per batch
	BatchSize int
	// micro batching time cycle
	BatchCycle time.Duration
	// a channel through which jobs are injected to this microbatcher
	inputChannel chan jobWrapper
	// wrapper for job to contain client channel, to notify of result
	jobItems map[string]jobWrapper
	// concurrency control for JobItems
	mutex *sync.RWMutex
	// micro-cycles ticker
	metronome *time.Ticker
	// we need to wait for all jobs done before returning to client.
	shutdownWait *sync.WaitGroup
	//responsilbe dispatching numbe of jobs to BatchProcessor
	dispatcher *dispatcher
}

//This metod is part of initialization ritual for MicroBatcher.
//Set up the batching timer, and channel signals.
func (mb *MicroBatcher) startBatching() {

	for {
		select {
		//The internal 	tempo controller of dispatcher
		case <-mb.metronome.C:
			mb.mutex.Lock()
			if len(mb.jobItems) != 0 {

				mb.mutex.Unlock()
				mb.shutdownWait.Add(1)
				mb.dispatchJobs()

			} else {
				mb.mutex.Unlock()
			}

		case newEndJob, ok := <-mb.inputChannel:
			//shut down has been requested
			if !ok {
				mb.shutdownWait.Add(1)
				mb.dispatchJobs()
				return
			}

			//Jobs will be dispacthed either by time or by filling the JobItems list.
			mb.mutex.Lock()
			mb.jobItems[newEndJob.theJob.ID] = newEndJob
			mb.mutex.Unlock()
			if len(mb.jobItems) == mb.BatchSize {
				mb.shutdownWait.Add(1)
				mb.dispatchJobs()
			}
		}
	}
}

//Get jobs dispatched.
func (mb *MicroBatcher) dispatchJobs() {
	mb.mutex.Lock()
	toBeDispatchted := make([]jobWrapper, 0)
	for k := range mb.jobItems {
		toBeDispatchted = append(toBeDispatchted, mb.jobItems[k])
		delete(mb.jobItems, k)
	}
	mb.mutex.Unlock()
	mb.dispatcher.dispatch(toBeDispatchted)
}

// Run Adds a job to microbatcher and waits for response in own groutine
func (mb *MicroBatcher) Run(job Job) (*JobResult, error) {

	jrChannel := make(chan JobResult)
	go func(chan JobResult) {
		//pass the job to batcher
		mb.inputChannel <- jobWrapper{theJob: job, responseChannel: jrChannel}
	}(jrChannel)

	defer close(jrChannel)
	//wait to get the response for the above mentioned job
	val := <-jrChannel

	return &val, nil
}

//Shutdown Call this method, and microbatcher knows its time to wrap up.
func (mb *MicroBatcher) Shutdown() {
	mb.metronome.Stop()
	close(mb.inputChannel)
	mb.shutdownWait.Wait()
}

//Start strats batching
func (mb *MicroBatcher) Start() {
	mb.metronome = time.NewTicker(mb.BatchCycle)
	go mb.startBatching()

}

//NewMicroBatcher instantiates and returns a new Moicro batcher
// params
// bp a BatchProcessor function.
// batchSize max size of each batch.
// batchCycle max time between dispatches
// returns a configured dispatcher
func NewMicroBatcher(bp BatchExecuteFn, batchSize int, batchCycle time.Duration) *MicroBatcher {
	waitSignal := &sync.WaitGroup{}

	dispatcher := &dispatcher{
		ProcessorFn: bp,
		flushWait:   waitSignal,
	}

	mbInstance := &MicroBatcher{
		BatchSize:    batchSize,
		BatchCycle:   batchCycle,
		inputChannel: make(chan jobWrapper),
		jobItems:     make(map[string]jobWrapper),
		mutex:        &sync.RWMutex{},
		dispatcher:   dispatcher,
		shutdownWait: waitSignal,
	}

	mbInstance.Start()
	return mbInstance
}
