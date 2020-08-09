## Micro-Batcher

Taking advantage of GO concurrency features to build a micro-batcher.


### A Bit Of Design

It seems , that micro batching can be broken into 3 parts, from functional perspective. As listed here :

* One part deals with client, to accept the process requests(, and possibly returns the results).
* Anotehr component monitors the incoming process requests, and timing, and triggers  requests the BatchProcessor to be invoked.
* The 3rd is a dispatcher(Worker), and invokes the BatchProcessor, and gets the reults back.

Also I need to mention I was between two ways choosing how the batcher should return results to a client.

* Get called and return the result at the same line of code.
* Get called and return results in a channel, as they arrive(, this one I *favour* the most).

I followed the first one as it seemed closer to what was stated to the requirement of the challenge was(, you know the interview sensitivity). 

### The Compenents
Considering the 3 functions named above here are the components, I came up with the following:

* Dispatcher : Accepts process requests from MicroBatcher, calls the BatchProcessor, and informs the client of results using a client specific channel.
* MicroBatcher: Responds to client process requests, and calls Dispatcher to process requests, when it decides so.

Obviously the first, and second functions are placed in Microbatcher.
### The Diagram

```console				

+--------+	        +-------------------------------+     		        +--------------+
| client +- Run(Job) -->+ IC(jobWrapper)=> Microbatcher |--- dispatch(jobs) --- + Dispatcher   +------- [BatchProcessor.Process]
+--------+  <---+       |	            		|	   	        |              | 
		|       +-------------------------------+		        +------+-------+
		+----------------------------------------------------------------------+	      								
```			      
### How To Use It

MicorBatcher is configurable, at initilisation time. Here are the confugurable values:
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


