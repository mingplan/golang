package main

import (
	"fmt"
	"os"
)

func main() {
	fileName := "./t.txt"
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for i := 0; i < 10; i++ {
		// n, err := f.Write([]byte(fmt.Sprintf("%d\n", i)))
		n, err := f.WriteString(fmt.Sprintf("%d\n", i))
		if err != nil {
			fmt.Println("write error:")
			fmt.Println(err.Error())
			os.Exit(2)
		}
		fmt.Println(fmt.Sprintf("writed %d bytes", n))
	}
}
