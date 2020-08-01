package microbatcher

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	processor := newSampleProceesor(time.Millisecond * 300)

	var wg sync.WaitGroup

	micro := NewMicroBatcher(processor.Execute, 2, time.Millisecond*1500)
	wg.Add(3)
	for _, i := range []int{1, 2, 3} {

		go func(number int) {

			jb := Job{Param: number, ID: fmt.Sprintf("ccddldwsw - %d", number)}
			result, err := micro.Run(jb)
			fmt.Printf("CMD .. main ..  CLIENT Got the result %v  ,  %v \n", result, err)
			wg.Done()
		}(i)
	}
	wg.Wait() //make sure all requests and resposnes are back before we are out
	fmt.Println("cmd Shutdown requested ... ")
	micro.Shutdown()
	fmt.Println("DONE for Now")

}
