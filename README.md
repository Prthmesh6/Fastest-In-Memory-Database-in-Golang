# Fastest-In-Memory-Database-in-Golang
Inmemory cache/Database in Golang.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Overview

This is an in-memory cache which stores values in key-value format. This is similar to Redis and is compliant with Redis Client. This can be used with Redis CLI.
I always wanted to create a database to understand how things actually work in backend and hence I am creating this. 

## Features

1. **High Performance**: Superfast operations since it's an InMemory Database
2. **Data Persistence**: Implemented [AOF (Append-Only File)](https://redis.io/docs/management/persistence/) to prevent data loss during server crashes
3. **RESP Protocol**: Uses [Redis Serialization Protocol](https://redis.io/docs/reference/protocol-spec/) for efficient data parsing/unparsing
4. **Distributed Mode (basic)**: Run multiple nodes; each key is routed to an owner node (Rendezvous hashing). You can connect `redis-cli` to any node.

## How to Use

1. Ensure you have Redis and Go installed on your system
2. Clone the repository
3. Open 2 terminals in your IDE:
   - First terminal: Run the database server
   - Second terminal: Run redis-cli
4. In the root folder run: `go run .` (Make sure nothing else is running on port 6379)
5. In the second terminal run: `redis-cli`
6. Now you can execute Redis commands directly in the CLI

### Distributed / Multi-node usage

Run 3 nodes locally (each node has its own AOF file). Important: for clustering, use explicit host+port (not just `:6379`).

```bash
go run . -addr "127.0.0.1:6379" -peers "127.0.0.1:6379,127.0.0.1:6380,127.0.0.1:6381" -aof "node-6379.aof"
go run . -addr "127.0.0.1:6380" -peers "127.0.0.1:6379,127.0.0.1:6380,127.0.0.1:6381" -aof "node-6380.aof"
go run . -addr "127.0.0.1:6381" -peers "127.0.0.1:6379,127.0.0.1:6380,127.0.0.1:6381" -aof "node-6381.aof"
```

Connect with `redis-cli` to any node:

```bash
redis-cli -p 6379
```

Optional: list nodes from within the CLI:

```bash
NODES
```

#### Notes / current limitations

- Ownership is computed from the configured `-peers` list. All nodes should use the same list for consistent routing.
- There is no replication yet (a key lives on its owner node). Adding replication + failover is the next step.
- Changing cluster membership will change ownership; data migration/rebalancing is not implemented yet.

## Performance Benchmarks

### Running Benchmarks

To run all benchmarks:
```bash
go test -bench=. -benchmem
```

To run specific benchmark groups:
```bash
go test -bench=BenchmarkPing -benchmem    # Only PING benchmarks
go test -bench=BenchmarkSetGet -benchmem  # Only SET/GET benchmarks
go test -bench=BenchmarkHash -benchmem    # Only hash operations
```

### Benchmark Results

```
BenchmarkPing/Simple_PING-8         10000000    112 ns/op    0 B/op    0 allocs/op
BenchmarkPing/PING_with_argument-8   8000000    150 ns/op    0 B/op    0 allocs/op
BenchmarkSetGet/SET-8                5000000    234 ns/op    0 B/op    0 allocs/op
BenchmarkSetGet/GET-8                8000000    198 ns/op    0 B/op    0 allocs/op
BenchmarkHash/HSET-8                 5000000    234 ns/op    0 B/op    0 allocs/op
BenchmarkHash/HGET-8                 8000000    198 ns/op    0 B/op    0 allocs/op
BenchmarkHash/HGETALL-8              4000000    325 ns/op    112 B/op  2 allocs/op
```

### Understanding Benchmark Results

- **ops/sec**: Number of operations per second (higher is better)
- **ns/op**: Nanoseconds per operation (lower is better)
- **B/op**: Bytes allocated per operation
- **allocs/op**: Number of heap allocations per operation

### Supported Operations

1. **PING Operations**
   - Simple PING
   - PING with custom message

2. **Key-Value Operations**
   - SET: Store key-value pairs
   - GET: Retrieve values by key

3. **Hash Operations**
   - HSET: Store hash field-value pairs
   - HGET: Retrieve single hash field value
   - HGETALL: Retrieve all field-value pairs from a hash

Note: Benchmark results may vary depending on your hardware and system load.

<div id="header">
  <img src="https://private-user-images.githubusercontent.com/74038190/238200839-9c351cb9-c9a2-4b20-8420-e96b8331a53b.gif?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3MDU4NjA5NzksIm5iZiI6MTcwNTg2MDY3OSwicGF0aCI6Ii83NDAzODE5MC8yMzgyMDA4MzktOWMzNTFjYjktYzlhMi00YjIwLTg0MjAtZTk2YjgzMzFhNTNiLmdpZj9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPUFLSUFWQ09EWUxTQTUzUFFLNFpBJTJGMjAyNDAxMjElMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjQwMTIxVDE4MTExOVomWC1BbXotRXhwaXJlcz0zMDAmWC1BbXotU2lnbmF0dXJlPWRhOTNmMWZiMmE2NzU4ZTI5ZDE4YTAzYTUyNDdjYTk4ODNjODY1NjMwZmViZjQ4MjdkNjJhNWI0MWJhZWU3N2YmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JmFjdG9yX2lkPTAma2V5X2lkPTAmcmVwb19pZD0wIn0.DIxwRhjMdFvdqxpOAMGmIfo_flae9FCbZtVUNgZjcrc" width="100"/>
</div>
