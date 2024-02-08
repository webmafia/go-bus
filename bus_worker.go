package bus

// Spawns a new worker.
func (b *Bus) Worker() {
	b.wg.Add(1)
	go b.worker()
}

func (b *Bus) worker() {
	defer b.wg.Done()

	for {
		ev, ok := <-b.queue

		if !ok {
			return
		}

		b.mu.RLock()
		for _, sub := range b.subs[ev.sub] {
			sub(ev.msg)
		}
		b.mu.RUnlock()
	}
}
