package main

import "sync"

type Job func()

type Pool struct{
	workerQueue			chan Job
	wg					sync.WaitGroup
}

//create a new pool

func NewPool(workerCount int) *Pool{
	pool := &Pool{
		workerQueue: make(chan Job),
	}

	pool.wg.Add(workerCount)

	for i := 0; i<workerCount; i++{
		go func ()  {
			defer pool.wg.Done()
			for job := range pool.workerQueue{
				job()
			}
		}()
	}
	return pool
}

func(p *Pool)Wait(){
	close(p.workerQueue)
	p.wg.Wait()
}

func(p *Pool)AddJob(job Job){
	p.workerQueue <- job
}



func main() {
	
	// Nothing here.
	// Run the following command in CLI:
	// go test -benchmem -bench BenchmarkConnections
}
