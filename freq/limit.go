package main

import (
	"sync"
	"sync/atomic"
)

type void struct{}

func NewGoroutinePool(size int) *GoroutinePool {
	return &GoroutinePool{
		channel: make(chan void, size),
	}
}

type GoroutinePool struct {
	channel chan void
	count   uint64
	wg      sync.WaitGroup
}

func (pool *GoroutinePool) Add() {
	atomic.AddUint64(&pool.count, 1)
	pool.wg.Add(1)
	pool.channel <- void{}
}

func (pool *GoroutinePool) Done() {
	<-pool.channel
	pool.wg.Done()
	atomic.AddUint64(&pool.count, ^uint64(0))
}

func (pool *GoroutinePool) Size() uint64 {
	return pool.count
}

func (pool *GoroutinePool) Wait() {
	pool.wg.Wait()
}
