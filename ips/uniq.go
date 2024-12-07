package ips

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

const maxCapacity int = 1024 * 1

var ErrCannotRead = fmt.Errorf("file cannot read")

type ips struct {
	sync.Mutex
	data map[string]int
}

func (i *ips) add(ip string) {
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
	return &ips{data: make(map[string]int, 0)}
}

func Uniq(filename string) (int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, ErrCannotRead
	}
	defer f.Close()

	var ips = newIPS()
	var wg sync.WaitGroup

	buf := make([]byte, maxCapacity)
	sc := bufio.NewScanner(f)
	sc.Buffer(buf, maxCapacity)

	for sc.Scan() {
		wg.Add(1)

		go func(ip string) {
			defer wg.Done()
			ips.add(ip)
		}(sc.Text())
	}
	wg.Wait()

	return ips.cnt(), nil
}
