package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func write(f *os.File, wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		f.WriteString(fmt.Sprintf("hello %d\n", i))
		time.Sleep(1 * time.Second)
	}
	wg.Done()
}

func read(f *os.File) {
	for {
		buffer := make([]byte, 100)
		n, err := f.Read(buffer)
		if err != nil {
			fmt.Println("read error")
			break
		}
		fmt.Println(string(buffer[0:n]))
	}
}

func main() {
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go write(w, wg)
	go read(r)
	wg.Wait()
}
