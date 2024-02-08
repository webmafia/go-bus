package bus

import "testing"

func BenchmarkBusPub(b *testing.B) {
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
