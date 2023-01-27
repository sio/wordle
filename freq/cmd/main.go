package main

import (
	"fmt"

	"github.com/sio/wordle"
	"github.com/sio/wordle/freq"
)

func main() {
	words := wordle.RussianWords()
	var chars freq.CharFreq

	chars.Update(words)

	fmt.Println(len(words))
	fmt.Println(chars)

	var test wordle.Word
	test.Parse("мания")
	fmt.Println(chars.Score(test))
}
