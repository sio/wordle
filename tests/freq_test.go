package tests

import (
	"testing"

	"github.com/sio/wordle"
	"github.com/sio/wordle/freq"
)

func TestFreq(t *testing.T) {
	words := wordle.RussianWords()

	var chars freq.CharFreq
	chars.Update(words)

	var total, f freq.Frequency
	for _, f = range chars {
		total += f
	}
	if total != 1 {
		t.Errorf("sum of character probabilities is not 1.0: %v", total)
	}
}
