#!/bin/bash

echo "=== Starting TinyProto Performance Tests ==="
echo ""

if ! nc -z localhost 8080 2>/dev/null; then
    echo "ERROR: Server not running on port 8080"
    echo "Start it with: go run cmd/protoserver/main.go"
    exit 1
fi

echo "Server is running"
echo ""

echo "=== Apache Bench Tests ==="
echo "Test 1: Basic throughput (10k requests, 100 concurrent)"
ab -n 10000 -c 100 -q http://localhost:8080/ | grep -E "Requests per second|Time per request|Failed requests"
echo ""

echo "Test 2: Higher concurrency (10k requests, 500 concurrent)"
ab -n 10000 -c 500 -q http://localhost:8080/ | grep -E "Requests per second|Time per request|Failed requests"
echo ""

echo "=== wrk Latency Distribution Tests ==="
echo "Test 1: Standard load (4 threads, 100 connections, 30s)"
wrk -t4 -c100 -d30s --latency http://localhost:8080/
echo ""

echo "Test 2: High concurrency (8 threads, 1000 connections, 30s)"
wrk -t8 -c1000 -d30s --latency http://localhost:8080/
echo ""

echo "=== Custom Concurrent Load Test ==="
go run cmd/loadtest/main.go
echo ""

echo "=== Performance Tests Complete ==="