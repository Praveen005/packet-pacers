package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)


var cpuprofile = flag.String("cpu-profile", "", "write cpu profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		// runtime.SetCPUProfileRate(500)   // default is 100, gives more refined result but at the cost of latency
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	// Nothing here.
	// Run the following command in CLI:
	// go test -benchmem -bench BenchmarkConnections
}
