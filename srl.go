package main

import (
	"fmt"
	"os"
)

func main() {
	ttl, threshold, err := readArgs()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	s := newServer(ttl, threshold)
	s.start()
}
