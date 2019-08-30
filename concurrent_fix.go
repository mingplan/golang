package main

import "fmt"
import "sync"
import "sync/atomic"

func main() {

	var count int32
	wg := &sync.WaitGroup{}
	process := func(i int, w *sync.WaitGroup) {
		for i := 0; i < 320; i++ {
			for {
				c := count
				if atomic.CompareAndSwapInt32(&count, c, c+1) {
					break
				}
			}
		}
		w.Done()
	}

	loops := 1000
	for i := 0; i < loops; i++ {
		wg.Add(1)
		go process(i, wg)
	}

	wg.Wait()

	fmt.Println("wait")
	fmt.Println(count)

}
