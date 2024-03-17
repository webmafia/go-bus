package bus

import (
	"sync/atomic"
)

type subscribers struct {
	subs    []atomic.Pointer[subscriber]
	buckets Topic
}

func newSubscribers(buckets int) *subscribers {
	return &subscribers{
		subs:    make([]atomic.Pointer[subscriber], buckets),
		buckets: Topic(buckets),
	}
}

type subscriber struct {
	next  atomic.Pointer[subscriber]
	topic Topic
}

func (s *subscribers) add(topic Topic) {
	sub := &subscriber{
		topic: topic,
	}

	target := &s.subs[topic%s.buckets]
	var ok bool

	for {
		ok = target.CompareAndSwap(nil, sub)

		if ok {
			break
		}

		target = &target.Load().next
	}
}

func (s *subscribers) count(topic Topic) (n int) {
	sub := s.subs[topic%s.buckets].Load()

	for sub != nil {
		if sub.topic == topic {
			n++
		}

		sub = sub.next.Load()
	}

	return
}
