package microbatcher

//Job jobs are processed by batch processor haing the following structure
type Job struct {
	//an imaginary math calculation input parameter
	Param int
	//Each job is identified by an identifier
	ID string
}

//JobResult processor will return results in this structure
type JobResult struct {
	//out come of the math processor
	Result int
	//Id of the corresponding Job
	JobID string
}

//JobWrapper an internally used strucutre.
//It holdsa channel to the clients. so other components can let the client know of resutls
type JobWrapper struct {
	theJob          Job
	responseChannel chan<- JobResult
}

// BatchExecuteFn the signature of the processor function
type BatchExecuteFn func(jobs []Job) []JobResult
