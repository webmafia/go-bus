package bus

import (
	"sync"

	"github.com/webmafia/fast"
)

type eventPool struct {
	pool sync.Pool
}

func (p *eventPool) acquire(size int) (e *event) {
	if v := p.pool.Get(); v != nil {
		e = v.(*event)
	} else {
		e = new(event)
	}

	if cap(e.msg) < size {
		e.msg = fast.MakeNoZero(size)
	} else {
		e.msg = e.msg[:size]
	}

	return
}

func (p *eventPool) release(e *event) {
	e.subKey.topic = 0
	e.subKey.tab = 0
	e.msg = e.msg[:0]

	p.pool.Put(e)
}
