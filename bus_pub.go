package bus

import "sync"

// Publish an event to a bus.
func Pub[T any](b *Bus, topic Topic, v *T) bool {
	return b.pub(topic, v, nil)
}

// Publish an event to a bus. After all subscribers have been handled, return the object to the provided pool.
func PubPool[T any](b *Bus, topic Topic, v *T, pool *sync.Pool) bool {
	return b.pub(topic, v, pool)
}

func (b *Bus) pub(topic Topic, msg any, pool *sync.Pool) bool {
	in := toIface(msg)

	b.queue <- event{
		sub: subscription{
			topic: topic,
			typ:   in.tab,
		},
		msg:  in.data,
		pool: pool,
	}

	return true
}
