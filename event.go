package bus

import "unsafe"

type event struct {
	msg unsafe.Pointer
	sub subscription
}
