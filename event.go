package bus

import (
	"sync"
	"unsafe"
)

type event struct {
	sub  subscription
	msg  unsafe.Pointer
	pool *sync.Pool
}
