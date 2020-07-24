package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/kavinho/microbatcher"
)

func main() {
	processor := microbatcher.NewSampleProceesor(time.Millisecond * 300)

	var wg sync.WaitGroup

	micro := microbatcher.NewMicroBatcher(processor.Execute, 2, time.Millisecond*1500)
	wg.Add(3)
	for _, i := range []int{1, 2, 3} {

		go func(number int) {

			jb := microbatcher.Job{Param: number, ID: fmt.Sprintf("ccddldwsw - %d", number)}
			result, err := micro.Run(jb)
			if err == nil {
				fmt.Printf("CMD .. main ..  CLIENT Got the result %v \n", result)
			} else {
				fmt.Printf("CMD .. main .. CLIENT Got ERROR %v \n", err)
			}
			wg.Done()
		}(i)
	}
	time.Sleep(time.Second * 1)
	fmt.Println("cmd Shutdown requested ... ")

	micro.Shutdown()
	fmt.Println("Waiting for all resutls to come back")
	wg.Wait()
	fmt.Println("DONE for Now")

}
