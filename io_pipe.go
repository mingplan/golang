package main

import (
	"fmt"
	"io"
	"sync"
	"time"
)

func write(w *io.PipeWriter, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		io.WriteString(w, fmt.Sprintf("test %d\n", i))
		time.Sleep(1 * time.Second)
	}
	w.Close()
	wg.Done()
}

func read(r *io.PipeReader) {
	for {
		buffer := make([]byte, 100)
		n, err := r.Read(buffer)
		if err == io.EOF {
			fmt.Println("EOF")
			fmt.Println(n)
			r.Close()
			break
		} else if err != nil {
			fmt.Println("read error:")
			fmt.Println(err.Error())
			break
		}
		fmt.Println("read content:")
		fmt.Println(string(buffer[0:n]))
	}
}

func main() {

	r, w := io.Pipe()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go write(w, wg)
	go read(r)
	wg.Wait()

}
