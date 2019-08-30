package main

import "fmt"
import "sync"

func main() {

	var count int
  wg := &sync.WaitGroup{}
	process := func(i int, w *sync.WaitGroup) {
		fmt.Println("procss")
    fmt.Println(i)
		for i := 0; i < 320; i++ {
      count++
		}
    w.Done()
	}

  // 增大loops的值，越容易出现并发导致count值不如预期
  loops := 1000
  for i := 0; i < loops; i++ {
    wg.Add(1)
    go process(i, wg)
  }

  wg.Wait()

	fmt.Println("wait")
	fmt.Println(count)

}
