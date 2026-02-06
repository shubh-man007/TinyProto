package response

import (
	"bytes"
	"strings"
	"testing"

	"github.com/shubh-man007/TinyProto/internal/headers"
)

func TestWriteStatusLineStandalone(t *testing.T) {
	var buf bytes.Buffer

	code, err := WriteStatusLine(&buf, StatusOK)
	if err != nil {
		t.Fatalf("WriteStatusLine returned error: %v", err)
	}
	if code != StatusOK {
		t.Fatalf("expected code %d, got %d", StatusOK, code)
	}
	if got := buf.String(); got != "HTTP/1.1 200 OK\r\n" {
		t.Fatalf("unexpected status line: %q", got)
	}

	buf.Reset()
	code, err = WriteStatusLine(&buf, StatusCode(999))
	if err == nil {
		t.Fatalf("expected error for unsupported status code")
	}
	if code != StatusUnrecog {
		t.Fatalf("expected StatusUnrecog, got %d", code)
	}
}

func TestGetDefaultHeadersAndWriteResHeaders(t *testing.T) {
	h := GetDefaultHeaders(4)

	if got := h.Get("Content-Length"); got != "4" {
		t.Fatalf("expected Content-Length 4, got %q", got)
	}
	if got := h.Get("Connection"); got != "close" {
		t.Fatalf("expected Connection close, got %q", got)
	}
	if got := h.Get("Content-Type"); got != "text/plain" {
		t.Fatalf("expected Content-Type text/plain, got %q", got)
	}

	var buf bytes.Buffer
	if err := WriteResHeaders(&buf, h); err != nil {
		t.Fatalf("WriteResHeaders returned error: %v", err)
	}

	out := buf.String()
	if !strings.HasSuffix(out, CRLF+CRLF) {
		t.Fatalf("expected final CRLF CRLF, got %q", out)
	}
	if !strings.Contains(out, "content-length: 4\r\n") {
		t.Fatalf("expected content-length header in %q", out)
	}
	if !strings.Contains(out, "connection: close\r\n") {
		t.Fatalf("expected connection header in %q", out)
	}
	if !strings.Contains(out, "content-type: text/plain\r\n") {
		t.Fatalf("expected content-type header in %q", out)
	}
}

func TestWriterStateMachineAndLogResponse(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	body := []byte("Hello")
	h := GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "text/plain")

	if err := w.WriteStatusLine(StatusOK); err != nil {
		t.Fatalf("WriteStatusLine returned error: %v", err)
	}
	if err := w.WriteHeaders(h); err != nil {
		t.Fatalf("WriteHeaders returned error: %v", err)
	}
	if _, err := w.WriteBody(body); err != nil {
		t.Fatalf("WriteBody returned error: %v", err)
	}
	if w.Status != WriterStatusDone {
		t.Fatalf("expected writer status %d, got %d", WriterStatusDone, w.Status)
	}

	logged := w.LogResponse(StatusOK, h, string(body))
	if !strings.HasPrefix(logged, "HTTP/1.1 200 OK\r\n") {
		t.Fatalf("unexpected logged status line: %q", logged)
	}
	if !strings.Contains(logged, "\r\n\r\nHello") {
		t.Fatalf("expected body in logged response, got %q", logged)
	}
}

func TestWriterWriteBodyAndTrailersStateErrors(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	if _, err := w.WriteBody([]byte("x")); err == nil {
		t.Fatalf("expected error when writing body before headers")
	}

	if err := w.WriteTrailers(headers.NewHeaders()); err == nil {
		t.Fatalf("expected error when writing trailers before body is done")
	}
}

func TestWriterWriteTrailersSuccess(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	body := []byte("OK")
	h := GetDefaultHeaders(len(body))

	if err := w.WriteStatusLine(StatusOK); err != nil {
		t.Fatalf("WriteStatusLine returned error: %v", err)
	}
	if err := w.WriteHeaders(h); err != nil {
		t.Fatalf("WriteHeaders returned error: %v", err)
	}
	if _, err := w.WriteBody(body); err != nil {
		t.Fatalf("WriteBody returned error: %v", err)
	}

	tr := headers.NewHeaders()
	tr.Set("Trailer", "X-Trace-ID")
	tr.Set("X-Trace-ID", "abc123")

	if err := w.WriteTrailers(tr); err != nil {
		t.Fatalf("WriteTrailers returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "X-Trace-ID: abc123\r\n") {
		t.Fatalf("expected trailer header in %q", out)
	}
	if !strings.HasSuffix(out, CRLF) {
		t.Fatalf("expected final CRLF at end of trailers, got %q", out)
	}
}

