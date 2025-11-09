# TinyProto

TinyProto is a minimal HTTP server implementation built from scratch using TCP in Go. The project demonstrates how HTTP works at the protocol level by implementing request parsing, header handling, and response generation from first principles.

## Overview

The project implements a basic HTTP/1.1 server with the following core features:

- TCP-level connection handling
- HTTP request parsing and validation 
- Header canonicalization and management
- Response construction utilities
- Support for request body handling based on Content-Length

## Project Structure

```
TinyProto/
├── cmd/
│   ├── nettest/           # Network testing utilities
│   │   ├── tcplistener/   # TCP connection diagnostic tool
│   │   └── udpsender/     # UDP testing utility
│   └── protoserver/       # Main HTTP server entry point
│       ├── assets/        
│       └── main.go
├── internal/
│   ├── headers/          # HTTP header parsing and management
│   ├── request/          # Request parsing and validation
│   ├── response/         # Response construction utilities
│   └── server/          # TCP server implementation
├── go.mod
└── go.sum
```

## Core Components

- **HTTP Parser**: Implements streaming request parsing with state machine
- **Header Management**: Case-insensitive header storage with duplicate handling
- **TCP Server**: Handles connection acceptance and concurrent request processing
- **Response Builder**: Utilities for constructing valid HTTP responses

## Getting Started

1. Run the HTTP server:
```bash
go run cmd/protoserver/main.go
```

2. Test with curl:
```bash
curl -v http://localhost:8080
```

3. For debugging request parsing:
```bash
go run cmd/tcplistener/main.go
```

## Implementation Details

The server implements HTTP/1.1 by:
1. Accepting TCP connections on a configured port
2. Parsing incoming bytes as HTTP requests
3. Validating request lines and headers
4. Processing request bodies based on Content-Length
5. Generating and sending HTTP responses

## Testing

The project includes unit tests for core components:
- Header parsing and management
- Request parsing and validation
- State machine transitions

Run tests:
```bash
go test ./internal/...
```
