package main

import (
	"fmt"
	"time"

	"github.com/joaofnds/chev"
)

type Sum struct {
	a, b int
}

func main() {
	chev.On(func(s Sum) { fmt.Printf("0: %d\n", s.a+s.b) })
	chev.Once(func(s Sum) { fmt.Printf("1: %d\n", s.a+s.b) })

	for i := 0; i < 5; i++ {
		chev.Send(Sum{a: i, b: 0})
	}

	<-time.After(time.Second)
}
