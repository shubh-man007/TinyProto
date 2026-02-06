package request

import (
	"bytes"
	"testing"
)

func BenchmarkRequestParsing(b *testing.B) {
	request := []byte("GET / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: bench\r\nAccept: */*\r\n\r\n")
	b.ReportAllocs()

	for b.Loop() {
		reader := bytes.NewReader(request)
		_, err := RequestFromReader(reader)
		if err != nil {
			b.Fatal(err)
		}
	}
}
