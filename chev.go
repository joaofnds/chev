package chev

import (
	"fmt"
	"sync"
)

var (
	lock = sync.RWMutex{}
	m    = make(map[string][]*Listener)
)

type Flag int

const (
	FlagReset Flag = 0
	FlagOnce  Flag = 1 << iota
)

type Handler[T any] func(t T)
type Listener struct {
	flags Flag
	event string
	f     Handler[any]
	in    chan any
}

func NewListener(event string, f Handler[any]) *Listener {
	l := &Listener{
		event: event,
		f:     f,
		flags: FlagReset,
		in:    make(chan any),
	}

	go func() {
		for t := range l.in {
			switch {
			case l.IsOnce():
				l.f(t)
				Off(l)
				return
			default:
				l.f(t)
			}
		}
	}()

	return l
}

func (l *Listener) Once() *Listener {
	l.flags |= FlagOnce
	return l
}

func (l *Listener) IsOnce() bool {
	return l.flags&FlagOnce != 0
}

func Off(l *Listener) {
	lock.Lock()
	defer lock.Unlock()
	for i, listener := range m[l.event] {
		if listener == l {
			m[l.event] = append(m[l.event][:i], m[l.event][i+1:]...)
			close(l.in)
			break
		}
	}
}

func On[T any](f Handler[T]) *Listener {
	name := fmt.Sprintf("%T", *new(T))
	h := NewListener(name, wrap(f))
	lock.Lock()
	defer lock.Unlock()
	m[name] = append(m[name], h)
	return h
}

func Once[T any](f Handler[T]) *Listener {
	name := fmt.Sprintf("%T", *new(T))
	h := NewListener(name, wrap(f)).Once()
	lock.Lock()
	defer lock.Unlock()
	m[name] = append(m[name], h)
	return h
}

func Send[T any](t T) {
	name := fmt.Sprintf("%T", t)

	lock.RLock()
	defer lock.RUnlock()

	for _, l := range m[name] {
		go func(l *Listener) {
			defer kalm()
			l.in <- t
		}(l)
	}
}

func wrap[T any](f Handler[T]) Handler[any] {
	return func(t any) { f(t.(T)) }
}

func kalm() { _ = recover() }
