package server

import (
	"bufio"
	"bytes"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shubh-man007/TinyProto/internal/request"
	"github.com/shubh-man007/TinyProto/internal/response"
)

func TestHandlerErrorWriteErrorResponse(t *testing.T) {
	var buf bytes.Buffer

	herr := &HandlerError{
		Code:    response.StatusBadRequest,
		Message: "bad request",
	}

	if err := herr.WriteErrorResponse(&buf); err != nil {
		t.Fatalf("WriteErrorResponse returned error: %v", err)
	}

	out := buf.String()
	if !strings.HasPrefix(out, "HTTP/1.1 400 Bad Request\r\n") {
		t.Fatalf("unexpected status line in error response: %q", out)
	}
	if !strings.HasSuffix(out, "bad request") {
		t.Fatalf("expected error message in body, got %q", out)
	}
}

func TestServeAndHandleOK(t *testing.T) {
	errCh := make(chan error, 1)

	handler := func(w *response.Writer, req *request.Request) {
		body := []byte("OK")
		h := response.GetDefaultHeaders(len(body))
		if err := w.WriteStatusLine(response.StatusOK); err != nil {
			errCh <- err
			return
		}
		if err := w.WriteHeaders(h); err != nil {
			errCh <- err
			return
		}
		if _, err := w.WriteBody(body); err != nil {
			errCh <- err
			return
		}
		close(errCh)
	}

	s, err := Serve(0, handler)
	if err != nil {
		t.Fatalf("Serve returned error: %v", err)
	}
	defer s.Close()

	addr, ok := s.listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("listener address is not *net.TCPAddr")
	}

	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(addr.Port))
	if err != nil {
		t.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	reqStr := "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"
	if _, err := conn.Write([]byte(reqStr)); err != nil {
		t.Fatalf("failed to write request: %v", err)
	}

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	reader := bufio.NewReader(conn)
	statusLine, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read status line: %v", err)
	}

	if !strings.HasPrefix(statusLine, "HTTP/1.1 200 OK") {
		t.Fatalf("unexpected status line: %q", statusLine)
	}

	if err, ok := <-errCh; ok && err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
}

