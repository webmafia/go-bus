package bus

import (
	"context"
	"testing"
)

func BenchmarkBusPubVal(b *testing.B) {
	type Foobar struct {
		Id int
	}

	eb := NewBus(context.Background(), 32)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		PubVal(eb, 0, &Foobar{Id: 1})
	}

	b.StopTimer()
	eb.Close()
}

func BenchmarkBusPubVal_Parallell(b *testing.B) {
	type Foobar struct {
		Id int
	}

	eb := NewBus(context.Background(), 32)
	v := &Foobar{Id: 1}

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			PubVal(eb, 0, v)
		}
	})

	b.StopTimer()
	eb.Close()
}

func BenchmarkBusPub(b *testing.B) {
	eb := NewBus(context.Background(), 32)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Pub(eb, 0)
	}

	b.StopTimer()
	eb.Close()
}

func BenchmarkBusPub_Parallell(b *testing.B) {
	eb := NewBus(context.Background(), 32)

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			Pub(eb, 0)
		}
	})

	b.StopTimer()
	eb.Close()
}
