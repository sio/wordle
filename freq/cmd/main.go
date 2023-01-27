package main

import (
	"fmt"
	"sort"

	"github.com/sio/wordle"
	"github.com/sio/wordle/freq"
)

func main() {
	words := wordle.RussianWords()

	var chars freq.CharFreq
	chars.Update(words)

	sort.SliceStable(words, func(i, j int) bool {
		return chars.Score(words[i]) > chars.Score(words[j])
	})

	const outputSize = 10

	fmt.Printf("Dictionary contains %d words\n\n", len(words))

	fmt.Println("Highest scoring single words:")
	for i := 0; i < len(words) && i < outputSize; i++ {
		fmt.Println(" ", words[i].String())
	}

	fmt.Printf("\nCharacter frequency table:\n%v\n\n", &chars)

	var input = [...]string{
		"мания",
		"укроп",
		"бювет",
	}
	test := make([]wordle.Word, len(input))
	for i := 0; i < len(test); i++ {
		test[i].Parse(input[i])
	}
	baseline := chars.Score(test...)
	fmt.Printf("Baseline for %v is %.1f%%\n\n", input, baseline*100)
}
