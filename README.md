## Micro-Batcher

Taking advantage of GO concurrency features to build a micro-batcher.


### A bit of Design

It seems , that micro batching can be broken into 3 parts, from functional perspective.
* Deals with client, to accept the process requests, and returns the results.
* Anotehr component manages the incoming process requests, and controls when requests need to be sent to the BatchProcessor
* The 3rd is a dispatcher(Worker), and calls the requests to the BatchProcessor, and gets the reults back.



* Dispatcher : Accepts process requests from MicroBatcher, calls the BatchProcessor, and informs the client of results.
* MicroBatcher: Responds to client process requests, and calls Dispatcher to process requests.

### How to use it
MicorBatcher is configurable, at initilisation time. As exaplained Below:
```
BatchSize int // How many items ber batch must be used to call BatchProcessor
BatchCycle time.Duration //defines the micro batching time cycle.
```

Here is a example of how to call the little guy:
```Go
  //Just an imaginary BatchProcessor
	processor := microbatcher.NewSampleProceesor(time.Millisecond * 300)

	var wg sync.WaitGroup

	micro := microbatcher.NewMicroBatcher(processor.Execute, 3, time.Millisecond*1500)
	wg.Add(5)
	for _, i := range []int{1, 2, 3,4,5} {

		go func(number int) {

			jb := microbatcher.Job{Param: number, ID: fmt.Sprintf("ccddldwsw - %d", number)}
			result, err := micro.Run(jb)
      fmt.Printf("Got results %v , %v \n", result,err)
			}
			wg.Done()
		}(i)
	}
  //To indicate finish off , and return
	micro.Shutdown()
	wg.Wait() //just making sure all jobs have been returned, and we have printed them

```


