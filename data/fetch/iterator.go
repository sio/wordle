package main

type fetchIterator interface {
	Close() error
	Next() bool
	Value() string
}
