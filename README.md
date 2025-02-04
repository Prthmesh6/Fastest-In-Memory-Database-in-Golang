### Fastest-In-Memory-Database-in-Golang
Inmemory cache/Database in Golang.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Overview

This is an in memory cache which stores values in key-value format. This is similar to redis and I have made it compliant with Redis Client. This can be used with Redis CLI.
I always wanted to create a database to understand how the things actually work in backend and hence I am creating this. 

## How to use ?

1) You should have redis and Go installed in your system, then clone the repository
2) Open 2 terminals in your IDE, One for running our database server other for redis-cli
3) On root folder run "go run ." Make sure nothins else running on port 6379
4) Use "redis-cli" command to run redis client
5) Now on Redis-cli you can execute the redis commands directly.

<div id="header">
  <img src="https://private-user-images.githubusercontent.com/74038190/238200839-9c351cb9-c9a2-4b20-8420-e96b8331a53b.gif?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3MDU4NjA5NzksIm5iZiI6MTcwNTg2MDY3OSwicGF0aCI6Ii83NDAzODE5MC8yMzgyMDA4MzktOWMzNTFjYjktYzlhMi00YjIwLTg0MjAtZTk2YjgzMzFhNTNiLmdpZj9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPUFLSUFWQ09EWUxTQTUzUFFLNFpBJTJGMjAyNDAxMjElMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjQwMTIxVDE4MTExOVomWC1BbXotRXhwaXJlcz0zMDAmWC1BbXotU2lnbmF0dXJlPWRhOTNmMWZiMmE2NzU4ZTI5ZDE4YTAzYTUyNDdjYTk4ODNjODY1NjMwZmViZjQ4MjdkNjJhNWI0MWJhZWU3N2YmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JmFjdG9yX2lkPTAma2V5X2lkPTAmcmVwb19pZD0wIn0.DIxwRhjMdFvdqxpOAMGmIfo_flae9FCbZtVUNgZjcrc" width="100"/>
</div>

1) This database works superfast since it is InMemory Database. As it is InMemory, server crash can cause data loss. But as I love to solve such problems I have added logic which avoid data loss, I have implemented [AOF](https://redis.io/docs/management/persistence/)
2) For Parsing UnParsing the data we haven't used JSON, we have used [RESP](https://redis.io/docs/reference/protocol-spec/) that is Redis serialization protocol specification you can read more about it on Redis Website, I have implemented it in code. 
3) To be continued...


# Redis-like Server Benchmarks

This section shows the performance benchmarks for different operations in our Redis-like server implementation.

## Running the Benchmarks

To run the benchmarks, use:

```bash
go test -bench=. -benchmem
```

## Benchmark Results

Below are sample benchmark results showing operations per second, time per operation, and memory usage:

``` 
BenchmarkPing/Simple_PING-8 10000000 112 ns/op 0 B/op 0 allocs/op
BenchmarkPing/PING_with_argument-8 8000000 150 ns/op 0 B/op 0 allocs/op
BenchmarkSetGet/SET-8 5000000 234 ns/op 0 B/op 0 allocs/op
BenchmarkSetGet/GET-8 8000000 198 ns/op 0 B/op 0 allocs/op
BenchmarkHash/HSET-8 5000000 234 ns/op 0 B/op 0 allocs/op
BenchmarkHash/HGET-8 8000000 198 ns/op 0 B/op 0 allocs/op
BenchmarkHash/HGETALL-8 4000000 325 ns/op 112 B/op 2 allocs/op

```

### Understanding the Results

- **ops/sec**: Higher numbers mean more operations per second (better performance)
- **ns/op**: Nanoseconds per operation (lower is better)
- **B/op**: Bytes allocated per operation
- **allocs/op**: Number of heap allocations per operation

### Benchmark Details

1. **PING Operations**
   - Simple PING without arguments
   - PING with custom message

2. **Key-Value Operations**
   - SET: Setting key-value pairs
   - GET: Retrieving values by key

3. **Hash Operations**
   - HSET: Setting hash field-value pairs
   - HGET: Getting single hash field value
   - HGETALL: Retrieving all field-value pairs from a hash

## Running Specific Benchmarks

You can run specific benchmark groups using:

bash
go test -bench=BenchmarkPing -benchmem # Only PING benchmarks
go test -bench=BenchmarkSetGet -benchmem # Only SET/GET benchmarks
go test -bench=BenchmarkHash -benchmem # Only hash operations

```

This README section provides:
1. Instructions for running benchmarks
2. Sample benchmark results
3. Explanation of what the numbers mean
4. Details about each benchmark type
5. Instructions for running specific benchmark groups

You should replace the sample benchmark numbers with your actual results when you run the benchmarks on your system. The actual performance numbers will vary depending on the hardware and system load.

```