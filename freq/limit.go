package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type void struct{}

func NewGoroutinePool(size int) *GoroutinePool {
	return &GoroutinePool{
		maxsize: uint32(size),
	}
}

type GoroutinePool struct {
	count   uint32
	maxsize uint32
	wg      sync.WaitGroup
}

func (pool *GoroutinePool) Add() error {
	if pool.Size() >= pool.maxsize {
		return fmt.Errorf("goroutine pool is currently full")
	}
	atomic.AddUint32(&pool.count, 1)
	pool.wg.Add(1)
	return nil
}

func (pool *GoroutinePool) Done() {
	pool.wg.Done()
	atomic.AddUint32(&pool.count, ^uint32(0))
}

func (pool *GoroutinePool) Size() uint32 {
	return atomic.LoadUint32(&pool.count)
}

func (pool *GoroutinePool) Wait() {
	pool.wg.Wait()
}
