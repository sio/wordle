package freq

import (
	"fmt"

	"github.com/sio/wordle"
)

type Frequency float32

type CharFreq map[rune]Frequency

func (cf *CharFreq) Score(word wordle.Word) (Frequency, error) {
	seen := make(map[rune]bool)
	var score, current Frequency
	var ok bool
	for _, char := range word {
		if seen[char] {
			continue
		}
		seen[char] = true
		current, ok = (*cf)[char]
		if !ok {
			return 0, fmt.Errorf("character score not available: %c", char)
		}
		score += current
	}
	return score, nil
}

func (cf *CharFreq) Update(words []wordle.Word) {
	*cf = make(CharFreq)

	var total Frequency
	var word wordle.Word
	var char rune
	for _, word = range words {
		for _, char = range word {
			(*cf)[char]++
			total++
		}
	}
	for char := range *cf {
		(*cf)[char] /= total
	}
}
