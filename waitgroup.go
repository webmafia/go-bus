package bus

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type waitgroup struct {
	state atomic.Uint64
}

func waitgroupCount(orig *sync.WaitGroup) int32 {
	wg := (*waitgroup)(unsafe.Pointer(orig))
	state := wg.state.Load()
	return int32(state >> 32)
}
