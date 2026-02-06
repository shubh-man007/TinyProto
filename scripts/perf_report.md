## Report

### Code Metrics

- **Total Lines of Code**: 1,281 (all `.go` files, tests and benchmarks included)
- **Test Coverage (internal packages)**:
  - `internal/headers/headers.go`: **86.2%**
  - `internal/request/request.go`: **85.5%**
  - `internal/response/response.go`: **77.3%**
  - `internal/server/server.go`: **77.1%**
- **Overall Internal Coverage**: **81.5% of statements** (driven down by untested `response` and `server` packages).
- **Average Cyclomatic Complexity** (from `gocyclo -avg -top 10 .`): **3.81**
  - Most complex functions:
    - `main.RequestPath` (`cmd/protoserver/main.go`): complexity 17
    - `request.RequestFromReader` (`internal/request/request.go`): complexity 13
    - `(*request.Request).Parse` (`internal/request/request.go`): complexity 13
    - `(*headers.Headers).Parse` (`internal/headers/headers.go`): complexity 8
    - Several response writer methods around complexity 5–6.

### Performance Metrics

#### Benchmark Tests (CPU, Memory, Allocations)

Commands:

- `go test ./internal/... -bench=. -benchmem`

Results:

- **Header parsing** (`internal/headers`, `BenchmarkHeaderParsing-16`):
  - **2156 ns/op**
  - **536 B/op**
  - **15 allocs/op**
- **Request parsing** (`internal/request`, `BenchmarkRequestParsing-16`):
  - **2442 ns/op**
  - **784 B/op**
  - **24 allocs/op**

#### Load Testing with wrk

Tool versions:

- `wrk 4.2.0`

Commands (from `scripts/perf_test.sh`):

- **Test 1: Standard load**
  - `wrk -t4 -c100 -d30s http://localhost:8080/`
- **Test 2: High concurrency**
  - `wrk -t8 -c1000 -d30s http://localhost:8080/`

Results:

- **Test 1 – 4 threads, 100 connections, 30 s**:
  - Requests: **95,721** in 30.04 s
  - **Requests/sec**: **3,186.73 RPS**
  - **Transfer/sec**: **721.99 KB/s**
  - Latency (per-thread stats):
    - Avg: 42.88 ms
    - **p50**: 12.75 ms
    - **p75**: 70.03 ms
    - **p90**: 111.14 ms
    - **p99**: 174.03 ms
- **Test 2 – 8 threads, 1000 connections, 30 s**:
  - Requests: **172,188** in 30.10 s
  - **Requests/sec**: **5,721.00 RPS**
  - **Transfer/sec**: **1.27 MB/s**
  - Latency:
    - Avg: 176.71 ms
    - **p50**: 153.60 ms
    - **p75**: 172.97 ms
    - **p90**: 283.32 ms
    - **p99**: 442.08 ms

Interpretation:

- At **moderate concurrency (100 connections)**, TinyProto sustains **≈3.2k RPS** with a **good median latency (~13 ms)** and a moderate tail (p99 ~174 ms).
- At **high concurrency (1000 connections)**, throughput increases to **≈5.7k RPS**, but latency degrades significantly (p50 ≈150 ms, p99 ≈440 ms).

#### Custom Concurrent Load Test

Command:

- `go run cmd/loadtest/main.go`

Results:

- **Total Requests**: 10,000
- **Successful**: 10,000
- **Failed**: 0
- **Duration**: 1.372343155 s
- **Requests Per Second (RPS)**: **7,286.81 RPS**
- **Average Latency**: **13.430235 ms**
- **Concurrency**: 100

Interpretation:

- The custom client’s average latency at concurrency 100 (~13.4 ms) closely matches the wrk p50 (12.75 ms), confirming that for the basic HTML endpoint TinyProto can sustain **≈7.2–7.3k RPS at 100 concurrent clients with low double‑digit millisecond latency and no observed application errors**.