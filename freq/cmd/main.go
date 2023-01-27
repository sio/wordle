package main

import (
	"fmt"
	"sort"

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

// Search for starting words that score better than a baseline
func search(words []wordle.Word, baseline freq.Frequency, batchSize int) [][]wordle.Word {
	return [][]wordle.Word{}
}
