package medium

import (
	"bufio"
	"context"
)

func convRowsToIPs(rows []string) []IP {
	var ips = make([]IP, 0)
	for _, row := range rows {
		ips = append(ips, stringToIP(row))
	}
	return ips
}

var rowsBatch = make([]string, 0)

var reader = func(ctx context.Context, sc *bufio.Scanner, bs BatchSize) <-chan []IP {
	out := make(chan []IP)

	go func() {
		defer close(out)
		for sc.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				rowsBatch = append(rowsBatch, sc.Text())
				if len(rowsBatch) == int(bs) {
					out <- convRowsToIPs(rowsBatch)
					rowsBatch = []string{}
				}
			}
		}
		if len(rowsBatch) > 0 {
			out <- convRowsToIPs(rowsBatch)
		}
	}()

	return out
}
