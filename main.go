package main

import (
	"fmt"
	"runtime"
	"time"

	"Aggregator/proxy"
)

func main() {
	runtime.GOMAXPROCS(0)

	s := proxy.NewHttpServer("localhost:10000")
	errC := s.Open()

	go func() {
		time.Sleep(5 * time.Second)
		s.Shutdown()
	}()

	select {
	case err, ok := <-errC:
		if ok {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Printf("close actively\n")
		}
	}
}
