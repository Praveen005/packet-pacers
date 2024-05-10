package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

var portIdx atomic.Int64
var readersCount = 100

const (
	defaultReadBufferSize  = 8 * 1024 * 1024
	defaultWriteBufferSize = 8 * 1024 * 1024
)

func newUDPSocket() (fd int, port int, err error) {
	// Create local udp socket on any random port
	fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		panic(err)
	}

	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		panic(err)
	}

	addr := [4]byte{127, 0, 0, 1}
	portBase := 5000
	for {
		port = int(portIdx.Add(1) + int64(portBase))
		err = syscall.Bind(fd, &syscall.SockaddrInet4{
			Port: port,
			Addr: addr,
		})
		if err == nil {
			break
		}
		port++
	}
	err = syscall.SetNonblock(fd, false)
	return
}

// DO NOT MODIFY THIS FUNCTION
func testInit(readersCount int, verbose bool) (ports []int, readChan chan []byte, closeChan chan struct{}, err error) {
	ports = make([]int, readersCount)
	portsChan := make(chan int, readersCount)

	readChan = make(chan []byte, readersCount)
	closeChan = make(chan struct{}, 1)

	// Create readersCount udp sockets to read
	wg := sync.WaitGroup{}
	for i := 0; i < readersCount; i++ {
		wg.Add(1)
		go func(threadId int) {
			fd, port, err := newUDPSocket()
			if err != nil {
				return
			}

			portsChan <- port

			buf := make([]byte, 1500)
			wg.Done()
			for {
				select {
				case <-closeChan:
					return
				default:
					n, _, err := syscall.Recvfrom(fd, buf, 0)
					if err != nil {
						// Close the socket
						syscall.Close(fd)
						return
					}
					if verbose {
						data := string(buf[:n])
						fmt.Println("threadId", threadId, "read: ", n, "bytes", "data: ", data)
					}
					readChan <- buf[:n]
				}
			}
		}(i)
	}

	wg.Wait()

	for i := 0; i < readersCount; i++ {
		ports[i] = <-portsChan
	}

	return
}

// DO NOT MODIFY THIS FUNCTION
func getTestMsg() []byte {
	// Generate a 1500 byte random message
	buf := make([]byte, 1500)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return buf
}

// DO NOT MODIFY THIS FUNCTION
func waitForReaders(readChan chan []byte, b *testing.B) {
	// Wait for reader to read
	for i := 0; i < readersCount; i++ {
		select {
		case <-readChan:
		case <-time.After(1 * time.Second):
			b.Fatal("timeout") // This should not happen
		}
	}
}

func BenchmarkConnections(b *testing.B) {
	b.Run("baseline", func(b *testing.B) {
		BenchmarkRawUDP(b)
	})

	b.Run("Sample", func(b *testing.B) {
		BenchmarkSample(b)
	})
}

func BenchmarkRawUDP(b *testing.B) {
	b.StopTimer()

	// Create a file for storing the CPU profile
    f, err := os.Create("cpu2.out")
    if err != nil {
        fmt.Println("Could not create profile file:", err)
        return
    }
    defer f.Close()

    // Start CPU profiling
    if err := pprof.StartCPUProfile(f); err != nil {
        fmt.Println("Could not start profiling:", err)
        return
    }
    defer pprof.StopCPUProfile() // Stop profiling at the end

	testPort := 40101
	// Create a udp network connection
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: testPort,
	})
	if err != nil {
		b.Fatal(err)
	}

	ports, readChan, closeChan, err := testInit(readersCount, false)
	if err != nil {
		b.Fatal(err)
	}
	_ = readChan

	writer := func() {
		for i := 0; i < readersCount; i++ {
			buf := getTestMsg()
			_, err := conn.WriteTo(buf, &net.UDPAddr{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: ports[i],
			})
			if err != nil {
				b.Fatal(err)
			}
		}

		// End of code that you are permitted to modify
		waitForReaders(readChan, b)
	}

	// Sequential test
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		writer()
	}
	b.StopTimer()

	close(closeChan)
}

// using sendTo function and pre-allocated Buffers
func BenchmarkSample(b *testing.B) {
	b.StopTimer()

	// Create a file for storing the CPU profile
    f, err := os.Create("cpu1.out")
    if err != nil {
        fmt.Println("Could not create profile file:", err)
        return
    }
    defer f.Close()

    // Start CPU profiling
    if err := pprof.StartCPUProfile(f); err != nil {
        fmt.Println("Could not start profiling:", err)
        return
    }
    defer pprof.StopCPUProfile() // Stop profiling at the end

	ports, readChan, closeChan, err := testInit(readersCount, false) // DO NOT EDIT THIS LINE
	if err != nil {
		b.Fatal(err)
	}
	_ = readChan

	fd, _, err := newUDPSocket()
	if err != nil {
		panic(err)
	}

	defer func() {
		syscall.Close(fd)
	}()

	remoteAddr := make([]syscall.SockaddrInet4, readersCount)
	var addr *net.UDPAddr
	var raddr *syscall.SockaddrInet4

	for i := 0; i < readersCount; i++ {
		addr = &net.UDPAddr{Port: ports[i], IP: net.IPv4(127, 0, 0, 1)}
		raddr = &syscall.SockaddrInet4{Port: addr.Port, Addr: [4]byte{addr.IP[12], addr.IP[13], addr.IP[14], addr.IP[15]}}
		remoteAddr[i] = *raddr
	}

	preAllocBuffers := make([][]byte, readersCount)

	for i := range preAllocBuffers {
		preAllocBuffers[i] = getTestMsg()
	}
	
	var wg sync.WaitGroup
	writer := func() {
		for i := 0; i < readersCount; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				err = syscall.Sendto(fd, preAllocBuffers[i], syscall.MSG_DONTWAIT, &remoteAddr[i])
				if err != nil {
					panic(err)
				}

			}(i)

		}
		wg.Wait()
		waitForReaders(readChan, b)
	}

	// Sequential test
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		writer()
	}
	b.StopTimer()

	close(closeChan)
}
