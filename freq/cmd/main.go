package main

import (
	"fmt"
	"sort"
	"strings"

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

type searchResult struct {
	words []wordle.Word
	score freq.Frequency
	dict  *dictionary
}

func (r *searchResult) Append(word wordle.Word) {
	r.words = append(r.words, word)
	r.score = r.dict.Score(r.words...)
}

func (r *searchResult) Clear() {
	r.words = r.words[:0]
	r.score = 0
}

func (r *searchResult) String() string {
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
func (d *dictionary) Search(batchSize int, baseline freq.Frequency) { //[]wordle.Word {
	result := &searchResult{
		words: make([]wordle.Word, 0, batchSize),
		dict:  d,
	}
	for i := 0; i < len(*d.words); i++ {
		word := (*d.words)[i]
		score := d.freq.Score(word)
		if result.score+score*freq.Frequency(batchSize-len(result.words)) < baseline {
			fmt.Println("short-circuit")
			break
		}
		result.Append(word)
		if len(result.words) < batchSize {
			continue
		}
		fmt.Println(result)
		if result.score > baseline {
			fmt.Println(result)
			return
		} else {
			result.Clear()
		}
	}
	//return [][]wordle.Word{}
}
