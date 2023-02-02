package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sio/wordle"
)

const outputSize = 10

func main() {
	var dict dictionary
	dict.Fill(wordle.RussianWords())

	words := *dict.words
	fmt.Printf("Dictionary contains %d words\n\n", len(words))
	fmt.Println("Highest scoring single words:")
	for i := 0; i < len(words) && i < outputSize; i++ {
		fmt.Println(" ", words[i].String())
	}

	fmt.Printf("\nCharacter frequency table:\n%v\n\n", dict.freq)

	// Establish a baseline for further comparison
	var base = []string{
		"мания",
		"укроп",
		"бювет",
	}
	baseline := dict.ScoreString(base...)
	fmt.Printf("Baseline for %v is %v\n\n", base, baseline)

	fmt.Println("Highest scoring start combinations:")
	for _, size := range []int{2, 3, 4} {
		fmt.Println("  Number of words:", size)
		for ws := range dict.SearchTopScore(size) {
			fmt.Println("   ", ws, dict.Score(ws...))
		}
	}
}

type dictionary struct {
	words *[]wordle.Word
	freq  *wordle.CharFreq
}

func (d *dictionary) Fill(words *[]wordle.Word) {
	d.words = words
	d.freq = &wordle.CharFreq{}
	d.freq.Update(words)
	sort.SliceStable(*words, func(i, j int) bool {
		return d.freq.Score((*words)[i]) > d.freq.Score((*words)[j])
	})
}

func (d *dictionary) ScoreString(words ...string) wordle.Frequency {
	data := make([]wordle.Word, len(words))
	for i := 0; i < len(words); i++ {
		data[i].Parse(words[i])
	}
	return d.freq.Score(data...)
}

func (d *dictionary) Score(words ...wordle.Word) wordle.Frequency {
	return d.freq.Score(words...)
}

type searchState struct {
	dict     *dictionary
	cursor   int
	baseline *wordle.Frequency
	size     int
	words    []wordle.Word
	score    wordle.Frequency
	results  chan<- searchState
	wg       *sync.WaitGroup
}

func (r *searchState) Append(word wordle.Word) searchState {
	next := searchState{
		dict:     r.dict,
		cursor:   r.cursor,
		baseline: r.baseline,
		size:     r.size,
		words:    make([]wordle.Word, len(r.words), r.size),
		results:  r.results,
		wg:       r.wg,
	}
	for index, w := range r.words {
		next.words[index] = w
	}
	next.words = append(next.words, word)
	next.score = next.dict.Score(next.words...)
	if next.score > *next.baseline {
		*next.baseline = next.score
	}
	return next
}

func (r *searchState) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	for _, word := range r.words {
		builder.WriteString(word.String())
		builder.WriteRune(' ')
	}
	builder.WriteString(fmt.Sprintf("%v]", r.score))
	return builder.String()
}

// Search all word combinations for starting words that score better than a baseline
func (d *dictionary) SearchFull(size int, baseline wordle.Frequency) { //[]wordle.Word {
	results := make(chan searchState)
	wg := &sync.WaitGroup{}
	go d.recursiveSearch(searchState{
		dict:     d,
		size:     size,
		baseline: &baseline,
		results:  results,
		wg:       wg,
	})
	go func() {
		for {
			r := <-results
			fmt.Println(&r)
		}
	}()
	time.Sleep(1)
	wg.Wait()
	close(results)
	fmt.Println("Done.")
}

func (d *dictionary) recursiveSearch(search searchState) {
	search.wg.Add(1)
	defer search.wg.Done()

	if len(search.words) == search.size {
		if search.score < *search.baseline {
			return
		}
		search.results <- search
		return
	}
	for ; search.cursor < len(*d.words); search.cursor++ {
		word := (*d.words)[search.cursor]
		score := d.freq.Score(word)
		if search.score+score*wordle.Frequency(search.size-len(search.words)) < *search.baseline {
			break
		}
		d.recursiveSearch(search.Append(word))
	}
}
