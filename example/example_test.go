package main

import (
	"sync"
	"testing"

	"github.com/webmafia/bus"
)

func BenchmarkPub(b *testing.B) {
	eb := bus.NewBus(b.N)
	v := &Foobar{Id: 1}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bus.Pub(eb, 0, v)
	}

	b.StopTimer()
	eb.Close()
}

func BenchmarkSignal(b *testing.B) {
	cond := sync.NewCond(&sync.Mutex{})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cond.Signal()
	}
}

func BenchmarkLock(b *testing.B) {
	var f foo
	f.cond = sync.Cond{L: &f.mu}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f.doSomething()
	}

	b.StopTimer()
}

type foo struct {
	mu   sync.Mutex
	cond sync.Cond
	i    int
}

func (f *foo) doSomething() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.i++
	f.cond.Signal()
	return true
}
