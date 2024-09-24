package main

import "sync"

type DLQueue struct {
	mu sync.Mutex
	m  map[string]struct{}
}

func NewDLQueue() *DLQueue {
	return &DLQueue{
		m: make(map[string]struct{}),
	}
}

func (cm *DLQueue) Add(key string) {
	cm.mu.Lock()
	cm.m[key] = struct{}{}
	cm.mu.Unlock()
}
