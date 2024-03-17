package bus

import (
	"context"
	"unsafe"
)

type subKey struct {
	topic Topic
	tab   uintptr
}

type sub struct {
	key subKey
	cb  unsafe.Pointer
	add bool
}

// Subscribe to a topic with value. The provided callback will be called whenever a matching event arrives.
func SubVal[T any](b *Bus, topic Topic, cb func(context.Context, *T)) {
	var v *T
	b.sub(topic, toTab(v), unsafe.Pointer(&cb), true)
}

func Sub(b *Bus, topic Topic, cb func(context.Context)) {
	b.sub(topic, 0, unsafe.Pointer(&cb), true)
}

func UnsubVal[T any](b *Bus, topic Topic, cb func(context.Context, *T)) {
	var v *T
	b.sub(topic, toTab(v), unsafe.Pointer(&cb), false)
}

func Unsub(b *Bus, topic Topic, cb func(context.Context)) {
	b.sub(topic, 0, unsafe.Pointer(&cb), false)
}

func (b *Bus) sub(topic Topic, tab uintptr, cb unsafe.Pointer, add bool) {
	b.subQueue <- sub{
		key: subKey{
			topic: topic,
			tab:   tab,
		},
		cb:  cb,
		add: add,
	}
}
