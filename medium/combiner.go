package medium

import (
	"context"
	"sync"
)

var combiner = func(ctx context.Context, inputs ...<-chan []IP) <-chan []IP {
	out := make(chan []IP)

	var wg sync.WaitGroup
	mux := func(a <-chan []IP) {
		defer wg.Done()

		for in := range a {
			select {
			case <-ctx.Done():
			case out <- in:
			}
		}
	}
	wg.Add(len(inputs))

	for _, in := range inputs {
		go mux(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
