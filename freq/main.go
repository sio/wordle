package main

import (
	"fmt"
	"sort"

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
