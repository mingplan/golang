package main

import "fmt"

func main() {
  c := make(chan int, 1)
  //var c chan int
  // close a nil chan will panic
  close(c)

  for i := 0; i < 32; i++ {
    select {
      case t := <-c:
        // when c is closed, it will always get value
        fmt.Println(t)
        fmt.Println("ccc")
      default:
        fmt.Println("default")
    }
  }

}
