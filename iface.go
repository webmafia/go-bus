package bus

import "unsafe"

type iface struct {
	tab  unsafe.Pointer
	data unsafe.Pointer
}

//go:inline
func toIface(v any) iface {
	return *(*iface)(unsafe.Pointer(&v))
}
