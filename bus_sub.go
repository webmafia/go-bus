package bus

import "unsafe"

func Sub[T any](b *Bus, topic Topic, cb func(*T)) {
	var v *T
	in := toIface(v)
	b.sub(topic, in.tab, unsafe.Pointer(&cb))
}

func (b *Bus) sub(topic Topic, typ unsafe.Pointer, cb unsafe.Pointer) {
	sub := subscription{
		topic: topic,
		typ:   typ,
	}

	b.mu.Lock()
	b.subs[sub] = append(b.subs[sub], *(*func(unsafe.Pointer))(cb))
	b.mu.Unlock()
}
