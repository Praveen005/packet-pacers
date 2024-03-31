# UDP Network Writer Challenge

## Introduction
This is a hiring challenge for the intern role for core media team at our company. The goal of this challenge is to write a UDP network writer that performs significantly better than the baseline implementation.

## Challenge Description
In the `benchs_test.go` file, you will find a benchmark test that measures the performance of the UDP network writer. Your task is to improve the performance of the writer so that it is at least 2 times faster than the baseline implementation.

## Getting Started
To run the benchmark test, use the following command:
```bash
go test -benchmem -bench BenchmarkConnections
```
----------------------------------------------------------------------------------------------------------


## My Submission
 ### Optimizations
 1. Made the `conn` variable global for `BenchmarkSample` function:

 Every time the `BenchmarkRawUDP` function runs, a new connection `conn` is made. It adverselsy affects the performance of the program. In the `BenchmarkSample` now the global variable `conn` that we created is re-used, resulting it significant performance boost as seen in the stats below.

 ```
 goos: linux
goarch: amd64
pkg: dyte.io/net-assignment
cpu: AMD Ryzen 5 5500U with Radeon Graphics         
BenchmarkConnections/baseline-12                     170           6606608 ns/op          185810 B/op        700 allocs/op
BenchmarkConnections/Sample-12                      1430            841476 ns/op          160012 B/op        300 allocs/op
PASS
ok      dyte.io/net-assignment  5.915s
```
> **Conclusion 1**: The `Sample` benchmark outperforms the `baseline` benchmark in all key metrics, indicating that it is more efficient in terms of throughput, execution speed, memory usage, and memory allocation. This suggests that the optimizations made in the `Sample` benchmark have led to significant performance improvements.
