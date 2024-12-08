package medium

import (
	"bytes"
)

var worker = func(ipsBatch <-chan []IP) <-chan []IP {
	out := make(chan []IP)

	go func() {
		defer close(out)
		var ips = make([]IP, 0)

		for data := range ipsBatch {
			for _, ip := range data {
				if !containsIP(ips, ip) {
					ips = append(ips, ip)
				}
			}

		}
		out <- ips
	}()

	return out
}

func containsIP(ips []IP, ip IP) bool {
	for _, v := range ips {
		if bytes.Compare(v, ip) == 0 {
			return true
		}
	}
	return false
}
