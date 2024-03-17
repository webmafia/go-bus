package bus

import (
	"unsafe"

	"github.com/webmafia/fast"
)

// Publish an event to a bus.
func Pub(b *Bus, topic Topic) {
	e := b.eventPool.acquire(0)
	e.subKey = subKey{
		topic: topic,
	}

	b.queue <- e
}

// Publish an event to a bus with a value.
func PubVal[T any](b *Bus, topic Topic, v *T) {
	size := int(unsafe.Sizeof(*v))
	e := b.eventPool.acquire(size)
	e.subKey = subKey{
		topic: topic,
		tab:   toTab(v),
	}
	copy(e.msg, fast.PointerToBytes(v, size))

	b.queue <- e
}
