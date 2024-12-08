package medium

import (
	"bufio"
	"context"
	"os"
	"sync"
)

// according to the article
// https://medium.com/@snassr/processing-large-files-in-go-golang-6ea87effbfe2

//
// data -> reader -> workers -> combiner -> result
//

type IP string
type BatchSize int

const (
	Batch1K   BatchSize = 1000
	Batch10K  BatchSize = 10_000
	Batch100K BatchSize = 100_000
	Batch1M   BatchSize = 1_000_000
)

func stringToIP(s string) IP {
	return IP(s)
	// return IP(net.ParseIP(s))
}

type ips struct {
	sync.Mutex
	data map[IP]int
}

func (i *ips) add(ip IP) {
	i.Lock()
	defer i.Unlock()
	i.data[ip]++
}

func (i *ips) cnt() int {
	i.Lock()
	defer i.Unlock()
	return len(i.data)
}

func newIPS() *ips {
	return &ips{data: make(map[IP]int, 0)}
}

type Options struct {
	NumWorker int
	BatchSize BatchSize
}

func Uniq(filepath string, opt Options) (int, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return process(f, opt.NumWorker, opt.BatchSize), nil
}

func process(f *os.File, numWorkers int, batchSize BatchSize) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// split the file into chunks
	r := reader(ctx, bufio.NewScanner(f), batchSize)

	// send the chunks to the workers
	var workers = make([]<-chan []IP, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = worker(r)
	}

	// combine the results
	var ipsUniq = newIPS()
	var wg sync.WaitGroup

	for ips := range combiner(ctx, workers...) {
		wg.Add(len(ips))
		for _, ip := range ips {
			go func(ip IP) {
				defer wg.Done()
				ipsUniq.add(ip)
			}(ip)
		}
	}
	wg.Wait()

	return ipsUniq.cnt()
}
