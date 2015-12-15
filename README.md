# golang-func-vs-rpc-vs-protobufhttp2

## Introduction

The main aim of this repo is to benchmark the performance of normal "in-memory" function call vs RPC calls vs Protobuf RPC with HTTP2.

## Running the benchmark

`go get -u github.com/common-benchmarks/golang-func-vs-rpc-vs-protobufhttp2`

`cd %GOPATH%/src/github.com/common-benchmarks/golang-func-vs-rpc-vs-protobufhttp2`

`go test --bench .`

## Benchmark results

The results on my machine. 2nd gen i7, 12 GB RAM.

```
BenchmarkProtobufRpcCall-8         10000            165825 ns/op
BenchmarkRpcCall-8                 10000            112365 ns/op
BenchmarkNormalFunction-8       2000000000          0.60 ns/op
```
