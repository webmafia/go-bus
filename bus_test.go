package bus

import "testing"

func BenchmarkBusPub(b *testing.B) {
	type Foobar struct {
		Id int
	}

	eb := NewBus(b.N)
	v := &Foobar{Id: 1}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Pub(eb, 0, v)
	}

	b.StopTimer()
	eb.Close()
}

func BenchmarkBusPubDirect(b *testing.B) {
	// type Foobar struct {
	// 	Id int
	// }

	eb := NewBus(b.N)
	// v := &Foobar{Id: 1}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		eb.mu.Lock()

		// for b.full() && !b.closed {
		// 	b.pubCond.Wait()
		// }

		// if b.closed {
		// 	return false
		// }

		// Put value last in queue
		// i := b.idx(b.size)
		eb.size++
		// _ = i

		// in := toIface(msg)
		// _ = in

		// b.queue[i] = event{
		// 	msg: in.data,
		// 	sub: subscription{
		// 		topic: topic,
		// 		typ:   in.tab,
		// 	},
		// }

		// Signal to one (1) waiting worker that there are events in queue
		eb.workCond.Signal()
		eb.mu.Unlock()
	}

	b.StopTimer()
	eb.Close()
}
