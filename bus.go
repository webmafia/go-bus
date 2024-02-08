package bus

import (
	"sync"
	"unsafe"
)

// Creates a new bus with one worker. This should be enought in most cases. If the worker falls behind,
// spawn more workers with `NewBusWithWorkers`.
func NewBus(capacity int) *Bus {
	return NewBusWithWorkers(capacity, 1)
}

// Creates a new bus with N workers. This should be used when there are many or heavy subscribers.
func NewBusWithWorkers(capacity int, workers int) *Bus {
	b := NewBusWithoutWorker(capacity)

	for i := 0; i < workers; i++ {
		go b.Worker()
	}

	return b
}

// Creates a new bus without any worker. Spawn your own workers using the `Worker` method. If the capacity
// is filled up without any worker, there will be deadlock. This should only be used in rare cases.
func NewBusWithoutWorker(capacity int) *Bus {
	b := &Bus{
		queue:   make([]event, capacity),
		subs:    make(map[subscription][]func(unsafe.Pointer)),
		maxSize: capacity,
	}

	b.pubCond = sync.Cond{L: &b.mu}
	b.workCond = sync.Cond{L: &b.mu}

	return b
}

func Pub[T any](b *Bus, topic Topic, v *T) bool {
	return b._Pub(topic, v)
}

func Sub[T any](b *Bus, topic Topic, cb func(*T)) {
	var v *T
	in := toIface(v)
	b._Sub(topic, in.tab, unsafe.Pointer(&cb))
}

type event struct {
	msg unsafe.Pointer
	sub subscription
}

type subscription struct {
	topic Topic
	typ   unsafe.Pointer
}

type SubFn[T any] func(v *T)

type Bus struct {
	queue    []event
	subs     map[subscription][]func(unsafe.Pointer)
	start    int
	size     int
	maxSize  int
	mu       sync.Mutex
	pubCond  sync.Cond
	workCond sync.Cond
	closing  bool
	closed   bool
}

func (b *Bus) Size() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.size
}

func (b *Bus) MaxSize() int {
	return b.maxSize
}

func (b *Bus) Empty() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.empty()
}

//go:inline
func (b *Bus) empty() bool {
	return b.size == 0
}

func (b *Bus) Full() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.full()
}

//go:inline
func (b *Bus) full() bool {
	return b.size == b.maxSize
}

func (b *Bus) Close(graceful ...bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.closing = true

	if len(graceful) > 0 && graceful[0] {
		for !b.empty() {
			b.pubCond.Wait()
		}
	}

	b.closed = true

	b.workCond.Broadcast()
	b.pubCond.Broadcast()
}

func (b *Bus) _Sub(topic Topic, typ unsafe.Pointer, cb unsafe.Pointer) {
	sub := subscription{
		topic: topic,
		typ:   typ,
	}

	b.mu.Lock()
	b.subs[sub] = append(b.subs[sub], *(*func(unsafe.Pointer))(cb))
	b.mu.Unlock()
}

// Spawns a new worker. This function is blocking, and should be run in a goroutine.
func (b *Bus) Worker() {
	for {
		b.mu.Lock()

		for b.empty() && !b.closed {
			b.workCond.Wait()
		}

		if b.closed {
			b.mu.Unlock()
			return
		}

		ev, ok := b.next()
		closing := b.closing

		// Unlock early, so that we don't block any publisher in vain - but only if not closing
		if !closing || !ok {
			b.mu.Unlock()
		}

		if !ok {
			continue
		}

		for _, sub := range b.subs[ev.sub] {
			sub(ev.msg)
		}

		if closing {
			b.mu.Unlock()
		}

		// Signal to all waiting publishers that there is capacity left
		b.pubCond.Broadcast()
	}
}

func (b *Bus) next() (ev event, ok bool) {
	if b.empty() {
		return
	}

	i := b.start
	b.start = b.wrap(b.start + 1)
	b.size--

	return b.queue[i], true
}

func (b *Bus) _Pub(topic Topic, msg any) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	for b.full() && !b.closed {
		b.pubCond.Wait()
	}

	if b.closed {
		return false
	}

	// Put value last in queue
	i := b.idx(b.size)
	b.size++

	in := toIface(msg)

	b.queue[i] = event{
		msg: in.data,
		sub: subscription{
			topic: topic,
			typ:   in.tab,
		},
	}

	// Signal to one (1) waiting worker that there are events in queue
	b.workCond.Signal()
	return true
}

//go:inline
func (b *Bus) idx(i int) int {
	return b.wrap(b.start + i)
}

//go:inline
func (b *Bus) wrap(i int) int {
	return (i + b.maxSize) % b.maxSize
}
