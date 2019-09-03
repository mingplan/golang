package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	fileName := "t.txt"
	f, err := os.Open(fileName)
	if os.IsNotExist(err) {
		fmt.Println("file not exist")
		os.Exit(1)
	} else if err != nil {
		fmt.Println("open file error:")
		fmt.Println(err.Error())
		os.Exit(2)
	}
	r := []byte{}
	for {
		b := make([]byte, 10)
		n, err := f.Read(b)
		if err == io.EOF {
			fmt.Println("EOF")
			break
		} else if err != nil {
			fmt.Println("read error:")
			fmt.Println(err.Error())
			os.Exit(3)
		}
		fmt.Println("read content:")
		r = append(r, b[0:n]...)
	}
	fmt.Println(string(r))
}
