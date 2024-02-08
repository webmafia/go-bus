package bus

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func BenchmarkReflectType(b *testing.B) {
	var e event

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = reflect.TypeOf(e)
	}
}

func BenchmarkReflectTypeCb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = doReflect()
	}
}

func BenchmarkReflectTypeCbNothing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = doNothing()
	}
}

func BenchmarkReflectTypeCbInterface(b *testing.B) {
	var e event

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = doReflectAny(&e)
	}
}

func BenchmarkReflectTypeCbInterfaceUnsafe(b *testing.B) {
	var e event

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = doReflectUnsafe(&e)
	}
}

//go:noinline
func doReflect() reflect.Type {
	return reflect.TypeOf(event{})
}

//go:noinline
func doNothing() uint64 {
	return 1
}

//go:noinline
func doReflectAny(v any) reflect.Type {
	return reflect.TypeOf(v).Elem()
}

//go:noinline
func doReflectUnsafe(v any) unsafe.Pointer {
	return (*iface)(unsafe.Pointer(&v)).tab
}

func Example() {
	// var fn SubFn[any]

	var v any = event{}
	var v2 any = event{}
	var v3 any = &event{}
	var v4 any = &event{}

	fmt.Printf("%+v\n", *(*iface)(unsafe.Pointer(&v)))
	fmt.Printf("%+v\n", *(*iface)(unsafe.Pointer(&v2)))
	fmt.Printf("%+v\n", *(*iface)(unsafe.Pointer(&v3)))
	fmt.Printf("%+v\n", *(*iface)(unsafe.Pointer(&v4)))

	// Output: qwe
}
