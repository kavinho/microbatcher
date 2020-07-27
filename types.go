package microbatcher

//Job are processed by BatchProcessor having the following structure
type Job struct {
	//can be of any type
	Param interface{}
	//job identifier
	ID string
}

//JobResult processor will return results in this structure
type JobResult struct {
	//out come of the math processor
	Result interface{}
	//Id of the corresponding Job
	JobID string
}

//JobWrapper an internally used strucutre.
//It holds a channel to the client. Dispatcher will use this to let the client know of resutls
type JobWrapper struct {
	theJob          Job
	responseChannel chan<- JobResult
}

// BatchExecuteFn the signature of the processor function.
// The processor need to follow this signature, as we need to which job matches which Result,
// using Job ID
type BatchExecuteFn func(jobs []Job) []JobResult
