package headers

import "testing"

func BenchmarkHeaderParsing(b *testing.B) {
	headerBlock := []byte("Content-Type: application/json\r\nHost: localhost\r\nUser-Agent: bench\r\n\r\n")
	b.ReportAllocs()

	for b.Loop() {
		h := NewHeaders()
		if _, _, err := h.Parse(headerBlock); err != nil {
			b.Fatal(err)
		}
	}
}
