package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joaofnds/chev"
)

func main() {
	defer chev.Close()

	go func() {
		for t := range time.Tick(time.Second) {
			chev.Send(t)
		}
	}()

	for i := 0; i < 5; i++ {
		go func(i int) {
			for t := range chev.Listen[time.Time]() {
				fmt.Printf("%d: %s\n", i, time.Since(t))
			}
		}(i)
	}

	println("waiting...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGABRT)
	<-sigChan
}
