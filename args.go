package main

import (
	"fmt"
	"os"
	"strconv"
)

func readArgs() (int64, uint, error) {
	if len(os.Args) != 3 {
		return 0, 0, fmt.Errorf("usage: %s [ttl] [threshold]\n", os.Args[0])
	}
	ttl, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("ttl value must be an integer, got %s\n", os.Args[1])
	}
	threshold, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("threshold value must be an integer, got %s\n", os.Args[2])
	}
	return ttl, uint(threshold), nil
}
