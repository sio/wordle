package tests

import (
	"testing"

	"github.com/sio/wordle"

	"fmt"
	"regexp"
)

func TestWord(t *testing.T) {
	inputs := []string{
		"hello",
		"алоха",
		"house",
		"домик",
	}
	for _, expected := range inputs {
		t.Run(expected, func(t *testing.T) {
			var w wordle.Word
			w.Parse(expected)
			got := w.String()
			if got != expected {
				t.Errorf("got: %s %v, expected: %s", got, w, expected)
			}
		})
	}
}

func TestRussian(t *testing.T) {
	words := wordle.RussianWords()
	valid := regexp.MustCompile(fmt.Sprintf(`^[ЁёА-я]{%d}$`, wordle.WordSize))
	var failures int
	for _, word := range words {
		if !valid.MatchString(word.String()) {
			t.Errorf("failed regex validation: %s %v", word.String(), word)
			failures++
			if failures > 20 {
				t.Fatal("too many failures")
			}
		}
	}
}
