package main

import (
	"fmt"
	"ips_uniq/ips"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ips_uniq <file>")
		os.Exit(1)
	}

	var count, err = ips.Uniq(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(fmt.Sprintf("ips uniq: %v", count))
}
