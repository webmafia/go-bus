package bus

import "unsafe"

//go:inline
func toTab(v any) uintptr {
	return *(*uintptr)(unsafe.Pointer(&v))
}
