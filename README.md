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
    BenchmarkConnections/baseline-12                     199           5404626 ns/op          185777 B/op        700 allocs/op
    BenchmarkConnections/Sample-12                      1969            614839 ns/op          160009 B/op        300 allocs/op
    PASS
    ok      dyte.io/net-assignment  5.033s
    ```

    ```
    > time ./main

    real    0m0.075s
    user    0m0.001s
    sys     0m0.016s
    ```
    > **Conclusion 1**: The `Sample` benchmark outperforms the `baseline` benchmark in all key metrics, indicating that it is more efficient in terms of throughput, execution speed, memory usage, and memory allocation. This suggests that the optimizations made in the `Sample` benchmark have led to significant performance improvements.



2. Use of `sync.Pool`

    using `sync.Pool` is beneficial in many ways:

    1. Reduced Allocation Overhead: If the payload is being allocated and deallocated frequently (e.g., every time the loop iterates), this can lead to significant overhead. By using `sync.Pool`, we can reuse the same payload buffer across iterations, avoiding the need to allocate and deallocate memory each time.

    2. Garbage Collection Efficiency: By reusing objects, we reduce the number of objects that need to be collected, which can improve the efficiency of garbage collection.

    3. Memory Efficiency: Reusing objects can also reduce the overall memory footprint of your application, as we're not constantly allocating and deallocating memory.


    **Result**: It didn't quite improve our Benchmark result :p

    ```
    goos: linux
    goarch: amd64
    pkg: dyte.io/net-assignment
    cpu: AMD Ryzen 5 5500U with Radeon Graphics         
    BenchmarkConnections/baseline-12                     124           9076033 ns/op          184629 B/op        700 allocs/op
    BenchmarkConnections/Sample-12                       122           9383940 ns/op          184272 B/op        700 allocs/op
    PASS
    ok      dyte.io/net-assignment  4.677s
    ```
