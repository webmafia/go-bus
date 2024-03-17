package bus

import (
	"context"
	"unsafe"
)

// Spawns a new worker in the background.
func (b *Bus) startWorker() {
	b.wg.Add(1)
	go b.worker()
}

func (b *Bus) worker() {
	defer b.wg.Done()

	for {
		select {

		case <-b.ctx.Done():
			return

		case sub, ok := <-b.subQueue:
			if !ok {
				return
			}

			if sub.add {

				// Add new subscription (if it doesn't already exist)
				if b.findSub(sub.key, sub.cb) < 0 {
					b.subs[sub.key] = append(b.subs[sub.key], sub.cb)
				}

			} else if idx := b.findSub(sub.key, sub.cb); idx >= 0 {

				// Remove subscription
				l := len(b.subs[sub.key]) - 1
				last := b.subs[sub.key][l]
				b.subs[sub.key][idx] = last
				b.subs[sub.key] = b.subs[sub.key][:l]
			}

		case ev, ok := <-b.queue:
			if !ok {
				return
			}

			if ev.subKey.tab > 0 {
				msgPtr := *(*unsafe.Pointer)(unsafe.Pointer(&ev.msg))

				for _, sub := range b.subs[ev.subKey] {
					cb := *(*func(context.Context, unsafe.Pointer))(sub)
					cb(b.ctx, msgPtr)
				}
			}

			for _, sub := range b.subs[subKey{topic: ev.subKey.topic}] {
				cb := *(*func(context.Context))(sub)
				cb(b.ctx)
			}

			b.eventPool.release(ev)
		}

	}
}

func (b *Bus) findSub(key subKey, cb unsafe.Pointer) int {
	for i := range b.subs[key] {
		if b.subs[key][i] == cb {
			return i
		}
	}

	return -1
}
