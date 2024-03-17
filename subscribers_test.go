package bus

import (
	"fmt"
	"testing"
)

func BenchmarkSubscribersAdd(b *testing.B) {
	subs := newSubscribers(16)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		subs.add(Topic(i))
	}
}

func BenchmarkSubscribersAddParallell(b *testing.B) {
	subs := newSubscribers(64)

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		var i int

		for p.Next() {
			subs.add(Topic(i))
			i = (i + 1) % b.N
		}
	})
}

func Example_subscribers() {
	subs := newSubscribers(64)

	subs.add(1)
	subs.add(1)
	subs.add(1)

	fmt.Println(subs.count(1))

	// Output: mjau
}
