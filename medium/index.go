package medium

import (
	"bufio"
	"context"
	"net"
	"os"
)

// according to the article
// https://medium.com/@snassr/processing-large-files-in-go-golang-6ea87effbfe2

//
// data -> reader -> workers -> combiner -> result
//

const BatchSize1K = 1000

type IP net.IP

func stringToIP(s string) IP {
	return IP(net.ParseIP(s))
}

func Uniq(filepath string) (int, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return process(f, 10, BatchSize1K), nil
}

func process(f *os.File, numWorkers int, batchSize int) int {
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
	var ipsUniq = make([]IP, 0)

	for ips := range combiner(ctx, workers...) {
		for _, ip := range ips {
			if !containsIP(ipsUniq, ip) {
				ipsUniq = append(ipsUniq, ip)
			}
		}
	}
	return len(ipsUniq)
}
