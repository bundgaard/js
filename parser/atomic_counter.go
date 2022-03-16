package parser

import "sync"

type AtomicCounter struct {
	sync.Mutex
	count int64
}

func (ac *AtomicCounter) Add() {
	ac.Lock()
	defer ac.Unlock()
	ac.count++
}

func (ac *AtomicCounter) Get() int64 {
	ac.Lock()
	defer ac.Unlock()
	return ac.count
}

var counter = AtomicCounter{count: 0}
