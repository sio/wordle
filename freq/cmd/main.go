package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sio/wordle"
	"github.com/sio/wordle/freq"
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
	fmt.Printf("Baseline for %v is %.1f%%\n\n", base, baseline*100)

	dict.Search(3, baseline)
}

type dictionary struct {
	words *[]wordle.Word
	freq  *freq.CharFreq
}

func (d *dictionary) Fill(words *[]wordle.Word) {
	d.words = words
	d.freq = &freq.CharFreq{}
	d.freq.Update(words)
	sort.SliceStable(*words, func(i, j int) bool {
		return d.freq.Score((*words)[i]) > d.freq.Score((*words)[j])
	})
}

func (d *dictionary) ScoreString(words ...string) freq.Frequency {
	data := make([]wordle.Word, len(words))
	for i := 0; i < len(words); i++ {
		data[i].Parse(words[i])
	}
	return d.freq.Score(data...)
}

func (d *dictionary) Score(words ...wordle.Word) freq.Frequency {
	return d.freq.Score(words...)
}

type searchState struct {
	dict     *dictionary
	cursor   int
	baseline *freq.Frequency
	size     int
	words    []wordle.Word
	score    freq.Frequency
	results  chan<- searchState
	wg       *sync.WaitGroup
}

func (r searchState) Append(word wordle.Word) searchState {
	if r.words == nil {
		r.words = make([]wordle.Word, 0, r.size)
	}
	r.words = append(r.words, word)
	r.score = r.dict.Score(r.words...)
	if r.score > *r.baseline {
		*r.baseline = r.score
	}
	return r
}

func (r *searchState) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	for _, word := range r.words {
		builder.WriteString(word.String())
		builder.WriteRune(' ')
	}
	builder.WriteString(fmt.Sprintf("%.1f]", r.score*100))
	return builder.String()
}

// Search for starting words that score better than a baseline
func (d *dictionary) Search(size int, baseline freq.Frequency) { //[]wordle.Word {
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
		if search.score+score*freq.Frequency(search.size-len(search.words)) < *search.baseline {
			break
		}
		d.recursiveSearch(search.Append(word))
	}
}
