package main

import (
	"runtime"
	"sort"
	"time"

	"github.com/sio/wordle"
)

// Search for starting words with top possible score
func (dict *dictionary) SearchTopScore(size int) chan []wordle.Word {
	chars := make([]rune, len(*dict.freq))
	var index int
	for char := range *dict.freq {
		chars[index] = char
		index++
	}
	sort.SliceStable(chars, func(i, j int) bool {
		return (*dict.freq)[chars[i]] > (*dict.freq)[chars[j]]
	})

	keep := make(map[rune]bool)
	var target wordle.Frequency
	for _, char := range chars[:size*wordle.WordSize] {
		keep[char] = true
		target += (*dict.freq)[char]
	}

	search := &topScoreSearchState{
		target: target,
		keep:   keep,
		dict:   dict,
		result: make([]wordle.Word, 0, size),
		pool:   NewGoroutinePool(runtime.NumCPU() + 1),
	}
	results := make(chan []wordle.Word)

	search.pool.Add() // register the first goroutine
	go recursiveTopScoreSearch(results, search, true)

	go func() {
		time.Sleep(1 * time.Second)
		search.pool.Wait()
		close(results)
	}()
	return results
}

func recursiveTopScoreSearch(results chan<- []wordle.Word, search *topScoreSearchState, goro bool) {
	if goro {
		defer search.pool.Done()
	}

	if len(search.result) == cap(search.result) {
		results <- search.result
		return
	}
	for ; search.cursor < len(*search.dict.words); search.cursor++ {
		word := (*search.dict.words)[search.cursor]
		if !search.Valid(word) {
			continue
		}
		wordScore := search.dict.Score(word)
		ceiling := wordScore*wordle.Frequency(cap(search.result)-len(search.result)) + search.score
		delta := search.target - ceiling
		const threshold = 1e-7
		if delta > threshold {
			continue
		}
		err := search.pool.Add()
		if err == nil { // spawn new worker goroutine
			go recursiveTopScoreSearch(results, search.Append(word), true)
		} else { // not allowed to add a goroutine, continuing in the current one
			recursiveTopScoreSearch(results, search.Append(word), false)
		}
	}
}

type topScoreSearchState struct {
	cursor int
	target wordle.Frequency
	score  wordle.Frequency
	dict   *dictionary
	keep   map[rune]bool
	result []wordle.Word
	pool   *GoroutinePool
}

func (s *topScoreSearchState) Valid(word wordle.Word) bool {
	for _, char := range word {
		if !s.keep[char] {
			return false
		}
		if s.Seen(char) {
			return false
		}
	}
	return true
}

func (s topScoreSearchState) Append(word wordle.Word) *topScoreSearchState {
	out := s
	out.result = make([]wordle.Word, len(s.result), cap(s.result))
	copy(out.result, s.result)
	out.result = append(out.result, word)
	out.score += out.dict.Score(word)
	out.cursor++
	return &out
}

func (s *topScoreSearchState) Seen(char rune) bool {
	for _, word := range s.result {
		if word.Contains(char) {
			return true
		}
	}
	return false
}
