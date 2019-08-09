package main

import (
	"sync"
)

type tItem struct {
	createdAt int64
	cnt       uint
}

type tCachedItems struct {
	*sync.RWMutex
	data map[uint64]*tItem
}
