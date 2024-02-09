package bus

import (
	"runtime"
	"sync"
	"testing"
)

func BenchmarkBusPubWithoutWorker(b *testing.B) {
	type Foobar struct {
		Id int
	}

	eb := NewBusWithoutWorker(b.N)
	v := &Foobar{Id: 1}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Pub(eb, 0, v)
	}

	b.StopTimer()
	eb.Close()
}

func BenchmarkBusPub(b *testing.B) {
	type Foobar struct {
		Id int
	}

	eb := NewBus(32)
	v := &Foobar{Id: 1}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Pub(eb, 0, v)
	}

	b.StopTimer()
	eb.Close()
}

func BenchmarkBusPubParallell(b *testing.B) {
	type Foobar struct {
		Id int
	}

	eb := NewBus(32)
	v := &Foobar{Id: 1}

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			Pub(eb, 0, v)
		}
	})

	b.StopTimer()
	eb.Close()
}

func BenchmarkBusWithMultipleWorkers(b *testing.B) {
	type Foobar struct {
		Id int
	}

	eb := NewBusWithWorkers(32, runtime.GOMAXPROCS(-1))
	v := &Foobar{Id: 1}

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			Pub(eb, 0, v)
		}
	})

	b.StopTimer()
	eb.Close()
}

func BenchmarkBusPubNew(b *testing.B) {
	type Foobar struct {
		Id int
	}

	eb := NewBus(32)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Pub(eb, 0, &Foobar{Id: 1})
	}

	b.StopTimer()
	eb.Close()
}

func BenchmarkBusPubPool(b *testing.B) {
	type Foobar struct {
		Id int
	}

	eb := NewBus(32)
	pool := sync.Pool{
		New: func() any {
			return new(Foobar)
		},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := pool.Get().(*Foobar)
		PubPool(eb, 0, v, &pool)
	}

	b.StopTimer()
	eb.Close()
}
