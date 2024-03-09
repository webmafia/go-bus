package bench

import (
	"os"
	"testing"

	"github.com/edsrzf/mmap-go"
)

func BenchmarkFileWrite(b *testing.B) {
	f, err := os.OpenFile("test.tmp", os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		b.Fatal(err)
	}

	if err := f.Truncate(4096); err != nil {
		b.Fatal(err)
	}

	buf, err := mmap.Map(f, mmap.RDWR, 0)

	if err != nil {
		b.Fatal(err)
	}

	val := []byte("........")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		copy(buf[i%len(buf):], val)
	}

	b.StopTimer()

	if err := buf.Unmap(); err != nil {
		b.Fatal(err)
	}

	if err := f.Close(); err != nil {
		b.Fatal(err)
	}
}
