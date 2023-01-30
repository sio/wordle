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
	delta := total - 1
	if delta > 1e-6 || delta < -1e-6 {
		t.Errorf("sum of character probabilities is not 1.0: %v", total)
	}

	inputs := []string{
		"hello",
		"алоха",
		"house",
		"домик",
	}
	var w wordle.Word
	for _, input := range inputs {
		w.Parse(input)
		single := chars.Score(w)
		double := chars.Score(w, w)
		triple := chars.Score(w, w, w)
		if single != double {
			t.Errorf("single word score (%v) does not match double word score (%v) for %q", single, double, w)
		}
		if single != triple {
			t.Errorf("single word score (%v) does not match triple word score (%v) for %q", single, triple, w)
		}
	}

	w.Parse("яяяяя")
	wordScore := chars.Score(w)
	charScore := chars['я']
	if wordScore != charScore {
		t.Errorf("repeated word score (%v) does not match single char score (%v)", wordScore, charScore)
	}
}
