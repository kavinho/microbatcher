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
	InputChannel chan JobWrapper
	// wrapper for job to contain client channel, to notify of result
	JobItems map[string]JobWrapper
	// concurrency control for JobItems
	mutex *sync.RWMutex
	// micro-cycles ticker
	metronome *time.Ticker
	// we need to wait for all jobs done before returning to client.
	shutdownWait *sync.WaitGroup
	//responsilbe dispatching numbe of jobs to BatchProcessor
	Dispatcher *Dispatcher
}

//This metod is part of initialization ritual for MicroBatcher.
//Set up the batching timer, and channel signals.
func (mb *MicroBatcher) startBatching() {

	for {
		select {
		//The internal 	tempo controller of dispatcher
		case <-mb.metronome.C:

			//fmt.Printf("mb ..startBatching .. Tick.. %v \n", time.Now())
			mb.mutex.Lock()
			if len(mb.JobItems) != 0 {
				mb.mutex.Unlock()
				mb.shutdownWait.Add(1)
				//fmt.Println("mb ..startBatching .. AddWait")
				mb.dispatchJobs()

			} else {
				mb.mutex.Unlock()
			}

		case newEndJob, ok := <-mb.InputChannel:

			if !ok {
				//fmt.Println("<-mb.InputChannel mb.shutdownWait.Add(1)")
				mb.shutdownWait.Add(1)
				//fmt.Println("ok := <-mb.InputChannel: Sutdown in progress")
				mb.dispatchJobs()
				return
			}

			//Jobs will be dispacthed either by time or by filling the JobItems list.
			mb.mutex.Lock()
			mb.JobItems[newEndJob.theJob.ID] = newEndJob
			mb.mutex.Unlock()
			if len(mb.JobItems) == mb.BatchSize {
				//fmt.Println("<-mb.InputChannel mb.shutdownWait.Add(1)")
				mb.shutdownWait.Add(1)
				mb.dispatchJobs()
			}
		}
	}
}

//Get jobs dispatched.
func (mb *MicroBatcher) dispatchJobs() {
	mb.mutex.Lock()
	toBeDispatchted := make([]JobWrapper, 0)
	for k := range mb.JobItems {
		toBeDispatchted = append(toBeDispatchted, mb.JobItems[k])
		delete(mb.JobItems, k)
	}
	mb.mutex.Unlock()
	mb.Dispatcher.Dispatch(toBeDispatchted)
}

// Run Adds a job to microbatcher and waits for response in own groutine
func (mb *MicroBatcher) Run(job Job) (*JobResult, error) {

	jrChannel := make(chan JobResult)
	go func(chan JobResult) {
		//fmt.Printf("mb .. Run .. push  %v \n", job.ID)
		mb.InputChannel <- JobWrapper{theJob: job, responseChannel: jrChannel}

	}(jrChannel)
	defer close(jrChannel)
	//fmt.Println("mb .. Run .. wating for client channel to respond")
	val := <-jrChannel
	//fmt.Println("> > Run () Got from Response channel")

	return &val, nil
}

//Shutdown Call this method, and microbatcher knows its time to wrap up.
func (mb *MicroBatcher) Shutdown() {
	mb.metronome.Stop()
	close(mb.InputChannel)
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

	dispatcher := &Dispatcher{
		ProcessorFn: bp,
		flushWait:   waitSignal,
	}

	mbInstance := &MicroBatcher{
		BatchSize:    batchSize,
		BatchCycle:   batchCycle,
		InputChannel: make(chan JobWrapper),
		JobItems:     make(map[string]JobWrapper),
		mutex:        &sync.RWMutex{},
		Dispatcher:   dispatcher,
		shutdownWait: waitSignal,
	}

	mbInstance.Start()
	return mbInstance
}
