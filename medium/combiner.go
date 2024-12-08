package medium

import (
	"context"
	"sync"
)

var combiner = func(ctx context.Context, inputs ...<-chan []IP) <-chan []IP {
	out := make(chan []IP)

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, in := range inputs {
		go func(in <-chan []IP) {
			defer wg.Done()

			for ips := range in {
				select {
				case <-ctx.Done():
				case out <- ips:
				}
			}
		}(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
