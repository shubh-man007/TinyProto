// NOTE:
// Curl sends the header first and then the body, because of a missing CRLF it believes that there are more headers coming in
// and keeps the connection open, buffering the request payload.
// And when we interrupt the curl request, the data present in the buffer is sent to the server and is printed out in the console.

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/shubh-man007/TinyProto/internal/request"
)

func handleConn(conn net.Conn) {
	defer func() {
		log.Printf("Connection closed: %s", conn.RemoteAddr())
		conn.Close()
	}()

	log.Printf("Accepted connection from %s", conn.RemoteAddr())

	r, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error parsing request from %s: %v", conn.RemoteAddr(), err)
		return
	}

	fmt.Printf("Request line:\n")
	fmt.Printf("- Method: %s\n", r.RequestLine.Method)
	fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
	fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

	fmt.Printf("Headers:\n")
	for key, val := range r.Header.Iter() {
		fmt.Printf("- %s: %s\n", key, val)
	}

	fmt.Printf("Body:\n")
	fmt.Printf("%s\n", string(r.Body))
}

func main() {
	const port = ":8080"
	const name = "uppertcp"
	log.SetPrefix(name + "\t")

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen at port %s: %s\n", port, err.Error())
	}
	defer listener.Close()

	log.Printf("Listening at %s", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConn(conn)
	}
}
