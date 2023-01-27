package freq

import (
	"github.com/sio/wordle"

	"fmt"
	"sort"
	"strings"
)

type Frequency float32

type CharFreq map[rune]Frequency

func (cf *CharFreq) Score(words ...wordle.Word) Frequency {
	seen := make(map[rune]bool)
	var score, current Frequency
	var ok bool
	for _, word := range words {
		for _, char := range word {
			if seen[char] {
				continue
			}
			seen[char] = true
			current, ok = (*cf)[char]
			if !ok {
				return 0
			}
			score += current
		}
	}
	return score
}

func (cf *CharFreq) Update(words *[]wordle.Word) {
	*cf = make(CharFreq)

	var total Frequency
	var word wordle.Word
	var char rune
	for _, word = range *words {
		for _, char = range word {
			(*cf)[char]++
			total++
		}
	}
	for char := range *cf {
		(*cf)[char] /= total
	}
}

func (cf *CharFreq) String() string {
	chars := make([]rune, len(*cf))
	var index int
	for char := range *cf {
		chars[index] = char
		index++
	}
	sort.SliceStable(chars, func(i, j int) bool {
		return (*cf)[chars[i]] > (*cf)[chars[j]]
	})
	output := make([]string, len(chars))
	for i := 0; i < len(chars); i++ {
		output[i] = fmt.Sprintf("%c:%.1f", chars[i], (*cf)[chars[i]]*100)
	}
	return fmt.Sprintf("[%s]", strings.Join(output, ", "))
}
