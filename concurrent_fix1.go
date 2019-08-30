package main

import "fmt"
import "sync"
import "runtime"

func main() {
  // 通过设置只有一个进程运行，也能解决并发问题
	runtime.GOMAXPROCS(1)
	var count int64
	wg := &sync.WaitGroup{}
	process := func(i int, w *sync.WaitGroup) {
		fmt.Println("procss")
		fmt.Println(i)
		for i := 0; i < 320; i++ {
			count++
		}
		w.Done()
	}
  
	loops := 10000
	for i := 0; i < loops; i++ {
		wg.Add(1)
		go process(i, wg)
	}

	wg.Wait()

	fmt.Println("wait")
	fmt.Println(count)

}
