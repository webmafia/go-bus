package bus

func Pub[T any](b *Bus, topic Topic, v *T) bool {
	return b.pub(topic, v)
}

func (b *Bus) pub(topic Topic, msg any) bool {
	in := toIface(msg)

	b.queue <- event{
		msg: in.data,
		sub: subscription{
			topic: topic,
			typ:   in.tab,
		},
	}

	return true
}
