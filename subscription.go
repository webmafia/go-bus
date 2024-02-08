package bus

import "unsafe"

type subscription struct {
	topic Topic
	typ   unsafe.Pointer
}
