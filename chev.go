package chev

import (
	"fmt"
	"sync"
)

var (
	lock      = sync.RWMutex{}
	listeners = make(map[string][]chan any)
)

func Listen[T any]() <-chan T {
	name := fmt.Sprintf("%T", *new(T))

	c := make(chan any)
	lock.Lock()
	listeners[name] = append(listeners[name], c)
	lock.Unlock()
	return wrap[T](c)
}

func Send[T any](t T) {
	name := fmt.Sprintf("%T", t)

	lock.RLock()
	defer lock.RUnlock()
	for _, c := range listeners[name] {
		go func(c chan any) { c <- t }(c)
	}
}

func Close() {
	for _, cs := range listeners {
		for _, c := range cs {
			close(c)
		}
	}
}

func wrap[T any](in chan any) chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for v := range in {
			out <- v.(T)
		}
	}()
	return out
}
