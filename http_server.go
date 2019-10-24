package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var port = flag.Int("port", 8081, "port of server")
	flag.Parse()
	fmt.Println(*port)
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		s := fmt.Sprintf("hello, world!\ncode: %d\n", 201)
		io.WriteString(w, s)
	})
	mux.HandleFunc("/hi", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "hi!\n")
	})
	server := http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			fmt.Println("shut down error")
			fmt.Println(err.Error())
		}
		close(idleConnsClosed)
	}()

	server.ListenAndServe()
	<-idleConnsClosed
}
